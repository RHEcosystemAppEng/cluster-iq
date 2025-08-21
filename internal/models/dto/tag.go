package dto

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

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
func NewTagDTOResponseList(Tags []TagDTOResponse) *TagDTOResponseList {
	response := TagDTOResponseList{Tags: Tags}

	// Count only set list length > 0
	if count := len(Tags); count > 0 {
		response.Count = count
	}

	return &response
}

// TODO comments
type TagDTORequest struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	InstanceID string `json:"instanceID"`
}

// TODO: comments
func (e TagDTORequest) ToInventoryTag() *inventory.Tag {
	return inventory.NewTag(
		e.Key,
		e.Value,
		e.InstanceID,
	)
}

// TODO comments
type TagDTORequestList struct {
	Count int             `json:"count,omitempty"`
	Tags  []TagDTORequest `json:"tag"`
}

// TODO comments
func NewTagDTORequestList(Tags []TagDTOResponse) *TagDTOResponseList {
	Request := TagDTOResponseList{Tags: Tags}

	// Count only set list length > 0
	if count := len(Tags); count > 0 {
		Request.Count = count
	}

	return &Request
}
