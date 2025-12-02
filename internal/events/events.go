// Package events provides functionality for audit logging and event tracking.
package events

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
)

type Event struct {
	// Unique identifier for the log entry.
	ID int64 `db:"id"`
	// Action performed (e.g., "cluster_stopped").
	Action actions.ActionOperation `db:"action"`
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
