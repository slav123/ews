package ews

import (
	"encoding/xml"
	"errors"
)

type FindItemRequest struct {
	XMLName             struct{}            `xml:"m:FindItem"`
	Traversal           string              `xml:"Traversal,attr"`
	ItemShape           ItemShape           `xml:"m:ItemShape"`
	IndexedPageItemView IndexedPageItemView `xml:"m:IndexedPageItemView"`
	ParentFolderIds     ParentFolderIds     `xml:"m:ParentFolderIds"`
	Restriction         Restriction         `xml:"m:Restriction"`
}

type ItemShape struct {
	BaseShape            BaseShape            `xml:"t:BaseShape"`
	AdditionalProperties AdditionalProperties `xml:"t:AdditionalProperties"`
}

type ParentFolderIds struct {
	DistinguishedFolderId DistinguishedFolderId `xml:"t:DistinguishedFolderId"`
}

type Restriction struct {
	IsEqualTo *IsEqualTo `xml:"t:IsEqualTo"`
	Contains  *Contains  `xml:"t:Contains"`
}

type IsEqualTo struct {
	BaseFiltering
	FieldURIOrConstant FieldURIOrConstant `xml:"t:FieldURIOrConstant"`
}

type Contains struct {
	BaseFiltering
	Constant              []Constant `xml:"t:Constant,omitempty"`
	ContainmentMode       string     `xml:"ContainmentMode,attr"`
	ContainmentComparison string     `xml:"ContainmentComparison,attr"`
}

type BaseFiltering struct {
	AdditionalProperties `xml:"t:FieldURI"`
}
type FieldURIOrConstant struct {
	Constant []Constant `xml:"t:Constant,omitempty"`
}

type Constant struct {
	Value string `xml:"Value,attr,omitempty"`
}

type findItemResponseEnvelop struct {
	XMLName struct{}             `xml:"Envelope"`
	Body    findItemResponseBody `xml:"Body"`
}
type findItemResponseBody struct {
	FindItemResponse FindItemResponse `xml:"FindItemResponse"`
}

type FindItemResponse struct {
	FindItemResponseMessages FindItemResponseMessages `xml:"ResponseMessages"`
}

type FindItemResponseMessages struct {
	FindItemResponseMessage FindItemResponseMessage `xml:"FindItemResponseMessage"`
}

type FindItemResponseMessage struct {
	Response
	RootFolder *RootFolder `xml:"RootFolder"`
}

type RootFolder struct {
	Items                   FindItems `xml:"Items"`
	IndexedPagingOffset     int       `xml:"IndexedPagingOffset,attr"`
	TotalItemsInView        int       `xml:"TotalItemsInView,attr"`
	IncludesLastItemInRange bool      `xml:"IncludesLastItemInRange,attr"`
}

type FindItems struct {
	Message []FindItemMessage `xml:"Message"`
}

type FindItemMessage struct {
	ItemId  ItemId `xml:"ItemId"`
	Subject string `xml:"Subject"`
}

// https://docs.microsoft.com/en-us/exchange/client-developer/web-service-reference/finditem-operation
func FindItem(c Client, r *FindItemRequest) (*FindItemResponse, error) {

	xmlBytes, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}

	bb, err := c.SendAndReceive(xmlBytes)
	if err != nil {
		return nil, err
	}

	var soapResp findItemResponseEnvelop
	err = xml.Unmarshal(bb, &soapResp)
	if err != nil {
		return nil, err
	}

	if soapResp.Body.FindItemResponse.FindItemResponseMessages.FindItemResponseMessage.ResponseClass != ResponseClassSuccess {
		return nil, errors.New(soapResp.Body.FindItemResponse.FindItemResponseMessages.FindItemResponseMessage.MessageText)
	}

	return &soapResp.Body.FindItemResponse, nil
}
