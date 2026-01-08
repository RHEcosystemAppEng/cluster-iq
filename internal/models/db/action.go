package db

import (
	"database/sql"

	"github.com/lib/pq"
)

// ActionDBResponse represents the database schema for action details,
// linking each field to a corresponding column in the database.
type ActionDBResponse struct {
	ID        string         `db:"id"`
	Type      string         `db:"type"`
	Time      sql.NullTime   `db:"time"`
	CronExp   sql.NullString `db:"cron_exp"`
	Operation string         `db:"operation"`
	Status    string         `db:"status"`
	Enabled   bool           `db:"enabled"`
	ClusterID string         `db:"cluster_id"`
	Region    string         `db:"region"`
	AccountID string         `db:"account_id"`
	Instances pq.StringArray `db:"instances"`
}
