package dto

import "github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"

// Tag represents the data transfer object for a tag.
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

// TODO comments
type TagDTOResponse struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	InstanceID string `json:"instanceID"`
}

// TODO comments
type TagDTOResponseList struct {
	Count int              `json:"count,omitempty"`
	Tags  []TagDTOResponse `json:"tag"`
}

// TODO comments
func NewTagDTOResponseList(Tags []TagDTOResponse) TagDTOResponseList {
	response := TagDTOResponseList{Tags: Tags}

	// Count only set list length > 0
	if count := len(Tags); count > 0 {
		response.Count = count
	}

	return response
}
