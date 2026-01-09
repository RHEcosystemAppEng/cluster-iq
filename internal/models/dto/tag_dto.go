package dto

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TagDTORequest represents the data needed to create or update a tag.
type TagDTORequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
} // @name TagRequest

// ToInventoryTag converts TagDTORequest to inventory.Tag.
// InstanceID is set to empty string as it will be assigned later
// when the tag is associated with an instance during creation.
func (t TagDTORequest) ToInventoryTag() *inventory.Tag {
	return inventory.NewTag(
		t.Key,
		t.Value,
		"",
	)
}

// ToInventoryTagList converts a slice of TagDTORequest to a slice of inventory.Tag.
func ToInventoryTagList(dtos []TagDTORequest) *[]inventory.Tag {
	tags := make([]inventory.Tag, len(dtos))
	for i, tag := range dtos {
		tags[i] = *tag.ToInventoryTag()
	}

	return &tags
}

// ToTagDTORequest converts inventory.Tag to TagDTORequest,
// omitting the InstanceID field as it's not part of the request DTO.
func ToTagDTORequest(tag inventory.Tag) *TagDTORequest {
	return &TagDTORequest{
		Key:   tag.Key,
		Value: tag.Value,
	}
}

// ToTagDTORequestList converts a slice of inventory.Tag to a slice of TagDTORequest.
func ToTagDTORequestList(tags []inventory.Tag) *[]TagDTORequest {
	tagList := make([]TagDTORequest, len(tags))
	for i, tag := range tags {
		tagList[i] = *ToTagDTORequest(tag)
	}

	return &tagList
}

// TagDTOResponse represents the data transfer object for a tag response,
// containing tag key-value pairs with optional instance association.
type TagDTOResponse struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	InstanceID string `json:"instanceId,omitempty"`
} // @name TagResponse
