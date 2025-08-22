package dto

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
