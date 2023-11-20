package ewsutil

import (
	"github.com/slav123/ews"
)

func GetAttachment(c ews.Client, attachmentId ews.AttachmentId) (*ews.GetAttachments, error) {

	req := &ews.GetAttachmentRequest{
		Items: []ews.AttachmentItemMessage{{AttachmentId: attachmentId}},
	}

	resp, err := ews.GetAttachment(c, req)

	if err != nil {
		return nil, err
	}

	return &resp.GetAttachmentResponseMessages.GetAttachmentResponseMessage.Attachments, nil
}
