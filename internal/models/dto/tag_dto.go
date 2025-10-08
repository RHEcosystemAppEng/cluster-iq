package dto

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// Tag represents the data transfer object for a tag.
type TagDTORequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TODO: comments
func (t TagDTORequest) ToInventoryTag() *inventory.Tag {
	return inventory.NewTag(
		t.Key,
		t.Value,
		"",
	)
}

func (t TagDTORequest) ToInventoryTagList(dtos []TagDTORequest) *[]inventory.Tag {
	tags := make([]inventory.Tag, len(dtos))
	for i, tag := range dtos {
		tags[i] = *tag.ToInventoryTag()
	}

	return &tags
}

// TODO comments
type TagDTOResponse struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	InstanceID string `json:"instanceID"`
}
