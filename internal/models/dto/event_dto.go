package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
)

// EventDTORequest represents the data needed to create or update an event.
type EventDTORequest struct {
	ID             int64     `json:"id"`
	Action         string    `json:"action"`
	ResourceID     string    `json:"resourceId"`
	ResourceType   string    `json:"resourceType"`
	EventTimestamp time.Time `json:"timestamp"`
	Result         string    `json:"result"`
	Severity       string    `json:"severity"`
	Requester      string    `json:"requester"`
	Description    *string   `json:"description,omitempty"`
} // @name EventRequest

func (e EventDTORequest) ToModelEvent() *events.Event {
	return &events.Event{
		ID:             e.ID,
		Action:         actions.ActionOperation(e.Action),
		EventTimestamp: e.EventTimestamp,
		Description:    e.Description,
		ResourceID:     e.ResourceID,
		ResourceType:   e.ResourceType,
		Result:         e.Result,
		Severity:       e.Severity,
		Requester:      e.Requester,
	}
} // @name EventResponse

// ClusterEvent represents a generic cluster event
type ClusterEventDTOResponse struct {
	ID             int64                `json:"id"`
	Action         string               `json:"action"`
	ResourceID     *string              `json:"resourceId"`
	ResourceType   string               `json:"resourceType"`
	EventTimestamp time.Time            `json:"timestamp"`
	Result         actions.ActionStatus `json:"result"`
	Severity       string               `json:"severity"`
	Requester      string               `json:"requester"`
	Description    *string              `json:"description,omitempty"`
	ScheduleID     *int64               `json:"scheduleId,omitempty"`
} // @name ClusterEventResponse

// SystemEvent represents a system-level event, extending a cluster event with account details.
type SystemEventDTOResponse struct {
	ClusterEventDTOResponse
	ResourceName string `json:"resourceName"`
	AccountID    string `json:"accountId"`
	AccountName  string `json:"accountName"`
	Provider     string `json:"provider"`
} // @name SystemEventResponse
