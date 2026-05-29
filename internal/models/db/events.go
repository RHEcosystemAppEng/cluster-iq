package db

import (
	"database/sql"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
)

// ClusterEventDBResponse represents the database schema for cluster event details,
// linking each field to a corresponding column in the database.
type ClusterEventDBResponse struct {
	ID             int64                `db:"id"`
	EventTimestamp time.Time            `db:"event_timestamp"`
	TriggeredBy    string               `db:"triggered_by"`
	Action         string               `db:"action"`
	ResourceID     *string              `db:"resource_id"`
	ResourceType   string               `db:"resource_type"`
	Result         actions.ActionStatus `db:"result"`
	Description    *string              `db:"description,omitempty"`
	Severity       string               `db:"severity"`
}

// SystemEventDBResponse represents the database schema for system event details,
// extending ClusterEventDBResponse with account and provider information.
type SystemEventDBResponse struct {
	ClusterEventDBResponse
	AccountID sql.NullString `db:"account_id"`
	Provider  sql.NullString `db:"provider"`
}
