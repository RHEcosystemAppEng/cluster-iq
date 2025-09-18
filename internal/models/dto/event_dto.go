package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/audit"
)

type ClusterEventDTORequest struct {
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

func (c ClusterEventDTORequest) ToModelAuditLog() *audit.AuditLog {
	events := audit.AuditLog{
		ID:             c.ID,
		ActionName:     actions.ActionOperation(c.ActionName),
		EventTimestamp: c.EventTimestamp,
		Description:    c.Description,
		ResourceID:     c.ResourceID,
		ResourceType:   c.ResourceType,
		Result:         c.Result,
		Severity:       c.Severity,
		TriggeredBy:    c.TriggeredBy,
	}
	return &events
}

// ClusterEventDTORequestList represents the API Request containing a list of accounts.
type ClusterEventDTORequestList struct {
	Events []ClusterEventDTORequest `json:"events"` // List of accounts.
}

func (c ClusterEventDTORequestList) ToModelAuditLogList() *[]audit.AuditLog {
	var events []audit.AuditLog

	for _, event := range c.Events {
		events = append(events, *event.ToModelAuditLog())
	}

	return &events
}

// ClusterEvent represents a generic cluster event
type ClusterEventDTOResponse struct {
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

// TODO: comments
// ClusterEventDTOResponseList represents the API response containing a list of accounts.
type ClusterEventDTOResponseList struct {
	Count  int                       `json:"count,omitempty"` // Number of accounts, omitted if empty.
	Events []ClusterEventDTOResponse `json:"events"`          // List of accounts.
}

// TODO: comments
// NewClusterEventDTOResponseList creates a new ClusterEventDTOResponseList instance.
// It ensures that an empty array is returned if the input account list is empty.
//
// Parameters:
// - accounts: A slice of inventory.Account.
//
// Returns:
// - A pointer to an ClusterEventDTOResponseList.
func NewClusterEventDTOResponseList(events []ClusterEventDTOResponse) *ClusterEventDTOResponseList {
	response := ClusterEventDTOResponseList{Events: events}

	// Count only set list length > 0
	if count := len(events); count > 0 {
		response.Count = count
	}

	return &response
}

// SystemEvent represents a system-level event, extending a cluster event with account details.
type SystemEventDTOResponse struct {
	ClusterEventDTOResponse
	AccountID string `json:"accountId"`
	Provider  string `json:"provider"`
}

type SystemEventDTOResponseList struct {
	Count  int                      `json:"count,omitempty"` // Number of accounts, omitted if empty.
	Events []SystemEventDTOResponse `json:"events"`          // List of accounts.
}

func NewSystemEventDTOResponseList(events []SystemEventDTOResponse) *SystemEventDTOResponseList {
	response := SystemEventDTOResponseList{Events: events}

	// Count only set list length > 0
	if count := len(events); count > 0 {
		response.Count = count
	}

	return &response
}
