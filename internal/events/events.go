package events

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"go.uber.org/zap"
)

// AuditEvent represents an action taken within the system.
// It provides key metadata such as the action performed, the resource involved,
// the result, severity, and the user who triggered the event.
type AuditEvent struct {
	// Unique identifier for the log entry.
	ID int64 `json:"id"`
	// Name of the action performed (e.g., "cluster_stopped").
	ActionName string `json:"action_name"`
	// UTC timestamp of when the action occurred.
	EventTimestamp time.Time `json:"event_timestamp"`
	// Optional reason for the action; can be nil.
	Reason *string `json:"reason,omitempty"`
	// ID of the affected resource (e.g., cluster_id, instance_id).
	ResourceID string `json:"resource_id"`
	// Type of resource affected (e.g., "cluster", "instance").
	ResourceType string `json:"resource_type"`
	// Outcome of the action (e.g., "success", "error").
	Result string `json:"result"`
	// Log severity level (e.g., "info", "warning", "error").
	Severity string `json:"severity"`
	// User or system entity responsible for the action.
	TriggeredBy string `json:"triggered_by"`
}

type SQLEventClient interface {
	AddEvent(event models.AuditLog) (int64, error)
	UpdateEventStatus(eventID int64, result string) error
}

type EventService struct {
	sqlClient SQLEventClient
}

type EventOptions struct {
	Action       string
	Reason       *string
	ResourceID   string
	ResourceType string
	Result       string
	Severity     string
	TriggeredBy  string
}

func NewEventService(sqlClient SQLEventClient) *EventService {
	return &EventService{
		sqlClient: sqlClient,
	}
}

func (e *EventService) LogEvent(opts EventOptions) (int64, error) {
	log := logger.NewLogger()
	event := models.AuditLog{
		TriggeredBy:    opts.TriggeredBy,
		ActionName:     opts.Action,
		ResourceID:     opts.ResourceID,
		ResourceType:   opts.ResourceType,
		Result:         opts.Result,
		Reason:         opts.Reason,
		Severity:       opts.Severity,
		EventTimestamp: time.Now(),
	}
	eventID, err := e.sqlClient.AddEvent(event)
	if err != nil {
		log.Error("Failed to log event", zap.Error(err))
		return 0, err
	}
	return eventID, nil
}

func (e *EventService) UpdateEventStatus(eventID int64, result string) error {
	log := logger.NewLogger()
	err := e.sqlClient.UpdateEventStatus(eventID, result)
	if err != nil {
		log.Error("Failed to update event status", zap.Int64("event_id", eventID), zap.Error(err))
		return err
	}
	return nil
}
