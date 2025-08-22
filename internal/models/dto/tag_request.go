package dto

import "github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"

// TODO comments
type TagDTORequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TODO: comments
func (e TagDTORequest) ToInventoryTag() *inventory.Tag {
	return inventory.NewTag(
		e.Key,
		e.Value,
		"",
	)
}

// TODO comments
type TagDTORequestList struct {
	Tags []TagDTORequest `json:"tag"`
}

func (t *TagDTORequestList) ToInventoryTagList() *[]inventory.Tag {
	var tags []inventory.Tag

	for _, tag := range t.Tags {
		tags = append(tags, *tag.ToInventoryTag())
	}

	return &tags
}
