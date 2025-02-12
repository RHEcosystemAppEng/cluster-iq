package events

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"go.uber.org/zap"
)

// Event states
const (
	ResultPending = "Pending"
	ResultSuccess = "Success"
	ResultFailed  = "Failed"
)

// Event severity levels
const (
	SeverityInfo    = "Info"
	SeverityError   = "Error"
	SeverityWarning = "Warning"
)

type SQLEventClient interface {
	AddEvent(event models.AuditLog) (int64, error)
	UpdateEventStatus(eventID int64, result string) error
}

type EventService struct {
	sqlClient SQLEventClient
	logger    *zap.Logger
}

type EventTracker struct {
	eventID int64
	service *EventService
	logger  *zap.Logger
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
