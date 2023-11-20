package ewsutil

import (
	"github.com/slav123/ews"
)

func FindEmail(c ews.Client, q string, isEqual bool) (*ews.RootFolder, error) {

	req := &ews.FindItemRequest{
		Traversal: "Shallow",
		ItemShape: ews.ItemShape{BaseShape: ews.BaseShapeIdOnly, AdditionalProperties: ews.AdditionalProperties{
			FieldURI: []ews.FieldURI{
				{FieldURI: "item:Subject"},
			},
		}},
		IndexedPageItemView: ews.IndexedPageItemView{
			MaxEntriesReturned: 5,
			Offset:             0,
			BasePoint:          ews.BasePointBeginning,
		},
		Restriction:     ews.Restriction{},
		ParentFolderIds: ews.ParentFolderIds{DistinguishedFolderId: ews.DistinguishedFolderId{Id: "inbox"}},
	}

	if isEqual {
		req.Restriction.IsEqualTo = &ews.IsEqualTo{BaseFiltering: ews.BaseFiltering{
			AdditionalProperties: ews.AdditionalProperties{
				FieldURI: []ews.FieldURI{
					{FieldURI: "item:Subject"},
				},
			},
		},
			FieldURIOrConstant: ews.FieldURIOrConstant{
				Constant: []ews.Constant{
					{Value: q},
				},
			}}
		req.Restriction.Contains = nil
	} else {
		req.Restriction.IsEqualTo = nil
		req.Restriction.Contains = &ews.Contains{BaseFiltering: ews.BaseFiltering{
			AdditionalProperties: ews.AdditionalProperties{
				FieldURI: []ews.FieldURI{
					{FieldURI: "item:Subject"},
				},
			},
		},
			Constant: []ews.Constant{
				{Value: q},
			},
			ContainmentMode:       "Substring",
			ContainmentComparison: "IgnoreCase",
		}
	}
	resp, err := ews.FindItem(c, req)

	if err != nil {
		return nil, err
	}

	return resp.FindItemResponseMessages.FindItemResponseMessage.RootFolder, nil
}
