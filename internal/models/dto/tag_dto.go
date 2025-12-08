package dto

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// Tag represents the data transfer object for a tag.
type TagDTORequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
} // @name TagRequest

// TODO: comments
func (t TagDTORequest) ToInventoryTag() *inventory.Tag {
	return inventory.NewTag(
		t.Key,
		t.Value,
		"",
	)
}

func ToInventoryTagList(dtos []TagDTORequest) *[]inventory.Tag {
	tags := make([]inventory.Tag, len(dtos))
	for i, tag := range dtos {
		tags[i] = *tag.ToInventoryTag()
	}

	return &tags
}

func ToTagDTORequest(tag inventory.Tag) *TagDTORequest {
	return &TagDTORequest{
		Key:   tag.Key,
		Value: tag.Value,
	}
}

func ToTagDTORequestList(tags []inventory.Tag) *[]TagDTORequest {
	tagList := make([]TagDTORequest, len(tags))
	for i, tag := range tags {
		tagList[i] = *ToTagDTORequest(tag)
	}

	return &tagList
}

// TODO comments
type TagDTOResponse struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	InstanceID string `json:"instanceId,omitempty"`
} // @name TagResponse
