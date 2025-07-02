package dto

// Tag model generic tags as a Key-Value object
type Tag struct {
	// Tag's key
	Key string `json:"key"`

	// Tag's Value
	Value string `json:"value"`

	// InstanceName reference
	InstanceID string `json:"instance_id"`
}
