package dto

import "time"

// ClusterEvent represents a generic cluster event
type ClusterEvent struct {
	ID             int64     `json:"id"`
	ActionName     string    `json:"action_name"`
	EventTimestamp time.Time `json:"event_timestamp"`
	Description    *string   `json:"description,omitempty"`
	ResourceID     string    `json:"resource_id"`
	ResourceType   string    `json:"resource_type"`
	Result         string    `json:"result"`
	Severity       string    `json:"severity"`
	TriggeredBy    string    `json:"triggered_by"`
}

// SystemEvent represents a system-level event, extending a cluster event with account details.
type SystemEvent struct {
	ClusterEvent
	AccountID string `json:"account_id"`
	Provider  string `json:"provider"`
}
