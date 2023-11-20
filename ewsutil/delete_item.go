package ewsutil

import (
	"github.com/slav123/ews"
)

func DeleteEmail(c ews.Client, itemId ews.ItemId, deleteType string) (bool, error) {

	req := &ews.DeleteItemRequest{
		DeleteType: deleteType,
		Items:      []ews.DeleteItemMessage{{ItemId: itemId}},
	}
	_, err := ews.DeleteItem(c, req)

	if err != nil {
		return false, err
	}

	return true, nil
}
