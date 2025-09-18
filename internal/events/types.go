package events

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
)

// TODO Duplicated with AuditLog
type AuditEvent struct {
	// Unique identifier for the log entry.
	ID int64 `json:"id"`
	// Name of the action performed (e.g., "cluster_stopped").
	ActionName actions.ActionOperation `json:"action_name"`
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
	// Provider is the name of the infrastructure provider (e.g., "AWS", "GCP", "Azure").
	Provider string `json:"provider"`
}
