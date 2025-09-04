package audit

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
)

// AuditLog represents an immutable record of an action taken within the system.
// It provides key metadata such as the action performed, the resource involved,
// the result, severity, and the user who triggered the event.
//
//revive:disable:exported
type AuditLog struct {
	// Unique identifier for the log entry.
	ID int64 `db:"id"`
	// Name of the action performed (e.g., "cluster_stopped").
	ActionName actions.ActionOperation `db:"action_name"`
	// UTC timestamp of when the action occurred.
	EventTimestamp time.Time `db:"event_timestamp"`
	// Optional description for the action; can be nil.
	Description *string `db:"description"`
	// ID of the affected resource (e.g., cluster_id, instance_id).
	ResourceID string `db:"resource_id"`
	// Type of resource affected (e.g., "cluster", "instance").
	ResourceType string `db:"resource_type"`
	// Outcome of the action (e.g., "success", "error").
	Result string `db:"result"`
	// Log severity level (e.g., "info", "warning", "error").
	Severity string `db:"severity"`
	// User or system entity responsible for the action.
	TriggeredBy string `db:"triggered_by"`
}

// SystemAuditLogs extends AuditLog with cloud provider metadata.
type SystemAuditLogs struct {
	// Base audit log data.
	AuditLog
	// Cloud provider account ID.
	AccountID string `db:"account_id"`
	// Cloud provider name (e.g., AWS, GCP).
	Provider string `db:"provider"`
}
