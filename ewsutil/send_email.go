package ewsutil

import "github.com/slav123/ews"

// SendEmail helper method to send Message
func SendEmail(c ews.Client, to []string, subject, body string, isHtml bool, attachments []ews.CreateFileAttachment) error {

	m := ews.Message{
		ItemClass: "IPM.Note",
		Subject:   subject,
		Body: ews.Body{
			Body: []byte(body),
		},
		Attachments: &ews.CreateAttachments{CreateFileAttachment: attachments},
		Sender: ews.OneMailbox{
			Mailbox: ews.Mailbox{
				EmailAddress: c.GetUsername(),
			},
		},
	}

	if isHtml {
		m.Body.BodyType = "HTML"
	} else {
		m.Body.BodyType = "Text"
	}

	mb := make([]ews.Mailbox, len(to))
	for i, addr := range to {
		mb[i].EmailAddress = addr
	}
	m.ToRecipients.Mailbox = append(m.ToRecipients.Mailbox, mb...)

	return ews.CreateMessageItem(c, m)
}
