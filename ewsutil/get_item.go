package ewsutil

import (
	"github.com/slav123/ews"
)

func GetEmail(c ews.Client, itemId ews.ItemId, includeMimeContent ews.BooleanType) (*ews.ItemMessage, error) {

	req := &ews.GetItemRequest{
		ItemShape: ews.GetItemShape{
			BaseShape:          ews.BaseShapeDefault,
			IncludeMimeContent: includeMimeContent,
		},
		Items: []ews.GetItemMessage{{ItemId: ews.ItemId{
			Id:        itemId.Id,
			ChangeKey: itemId.ChangeKey,
		}}},
	}

	resp, err := ews.GetItem(c, req)

	if err != nil {
		return nil, err
	}

	return &resp.GetItemResponseMessages.GetItemResponseMessage.Items.Message[0], nil
}
