package dto

import "time"

// ClusterEvent represents a generic cluster event
type ClusterEvent struct {
	ID             int64     `json:"id"`
	ActionName     string    `json:"actionName"`
	EventTimestamp time.Time `json:"timestamp"`
	Description    *string   `json:"description,omitempty"`
	ResourceID     string    `json:"resourceId"`
	ResourceType   string    `json:"resourceType"`
	Result         string    `json:"result"`
	Severity       string    `json:"severity"`
	TriggeredBy    string    `json:"triggeredBy"`
}

// SystemEvent represents a system-level event, extending a cluster event with account details.
type SystemEvent struct {
	ClusterEvent
	AccountID string `json:"accountId"`
	Provider  string `json:"provider"`
}
