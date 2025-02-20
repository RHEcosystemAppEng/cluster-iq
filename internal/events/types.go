package events

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"go.uber.org/zap"
)

// Event states
const (
	ResultSuccess = "Success"
	ResultFailed  = "Failed"
	ResultPending = "Pending"
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
	Description  *string
	ResourceID   string
	ResourceType string
	Result       string
	Severity     string
	TriggeredBy  string
}

type AuditEvent struct {
	// Unique identifier for the log entry.
	ID int64 `json:"id"`
	// Name of the action performed (e.g., "cluster_stopped").
	ActionName string `json:"action_name"`
	// UTC timestamp of when the action occurred.
	EventTimestamp time.Time `json:"event_timestamp"`
	// Optional description for the action; can be nil.
	Description *string `json:"description,omitempty"`
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

// SystemAuditEvent extends AuditEvent with system-specific fields.
type SystemAuditEvent struct {
	AuditEvent
	// AccountID is the unique identifier of the cloud provider account.
	AccountID string `json:"account_id"`
	// Provider is the name of the cloud provider (e.g., "AWS", "GCP", "Azure").
	Provider string `json:"provider"`
}
