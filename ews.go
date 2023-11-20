package ews

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/Azure/go-ntlmssp"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
)

const (
	soapStart = `<?xml version="1.0" encoding="utf-8" ?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:m="http://schemas.microsoft.com/exchange/services/2006/messages" xmlns:t="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
<soap:Body>
`
	// Exchange2013_SP1
	soapEnd = `
</soap:Body>
</soap:Envelope>`
)

type Config struct {
	Dump    bool
	NTLM    bool
	SkipTLS bool
	Version string
	RTMutex *sync.Mutex
}

type Client interface {
	SendAndReceive(body []byte) ([]byte, error)
	GetEWSAddr() string
	GetUsername() string
}

type client struct {
	EWSAddr string
	Email   string
	Login   LoginStrategy
	config  *Config
}

func (c *client) GetEWSAddr() string {
	return c.EWSAddr
}

func (c *client) GetUsername() string {
	return c.Email
}

func NewClient(ewsAddr, username, password string, config *Config) Client {
	return NewClientWithLoginStrategy(ewsAddr, username, PlainLogin{
		Username: username, Password: password}, config)
}

func NewClientWithLoginStrategy(ewsAddr, email string, loginStrategy LoginStrategy, config *Config) Client {
	return &client{EWSAddr: ewsAddr, Email: email, Login: loginStrategy, config: config}
}

func (c *client) SendAndReceive(body []byte) ([]byte, error) {

	bb := []byte(soapStart)
	bb = append(bb, body...)
	bb = append(bb, soapEnd...)

	req, err := http.NewRequest("POST", c.EWSAddr, bytes.NewReader(bb))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(req.Body)
	logRequest(c, req)

	c.Login.SetLoginHeaders(req)
	req.Header.Set("Content-Type", "text/xml")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	applyConfig(c.config, client)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	logResponse(c, resp)

	if resp.StatusCode != http.StatusOK {
		return nil, NewError(resp)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, err
}

type mutexRT struct {
	mtx *sync.Mutex
	rt  http.RoundTripper
}

func (m *mutexRT) RoundTrip(r *http.Request) (*http.Response, error) {
	if m.mtx != nil {
		m.mtx.Lock()
		defer m.mtx.Unlock()
	}
	return m.rt.RoundTrip(r)
}

func applyConfig(config *Config, client *http.Client) {
	godebug := os.Getenv("GODEBUG")
	if !strings.Contains(godebug, "http2client=0") {
		if godebug != "" {
			godebug += ","
		}
		godebug += "http2client=0"
		_ = os.Setenv("GODEBUG", godebug)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if config.SkipTLS {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if config.NTLM {
		transport.MaxConnsPerHost = 1
		client.Transport = &mutexRT{mtx: config.RTMutex, rt: ntlmssp.Negotiator{
			RoundTripper: transport,
		}}
	} else {
		client.Transport = transport
	}
}

func logRequest(c *client, req *http.Request) {
	if c.config != nil && c.config.Dump {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Request:\n%v\n----\n", string(dump))
	}
}

func logResponse(c *client, resp *http.Response) {
	if c.config != nil && c.config.Dump {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Response:\n%v\n----\n", string(dump))
	}
}
