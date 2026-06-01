package db

import (
	"database/sql"

	"github.com/lib/pq"
)

// ActionDBResponse represents the database schema for action details,
// linking each field to a corresponding column in the database.
type ActionDBResponse struct {
	ID               string         `db:"id"`
	Type             string         `db:"type"`
	Time             sql.NullTime   `db:"time"`
	CronExp          sql.NullString `db:"cron_exp"`
	Operation        string         `db:"operation"`
	Status           string         `db:"status"`
	Enabled          bool           `db:"enabled"`
	TargetType       string         `db:"target_type"`
	SelectAll          bool           `db:"select_all"`
	ClusterID          sql.NullString `db:"cluster_id"`
	ClusterName        sql.NullString `db:"cluster_name"`
	Region             sql.NullString `db:"region"`
	AccountID          string         `db:"account_id"`
	Instances          pq.StringArray `db:"instances"`
	TargetAccountIDs   pq.StringArray `db:"target_account_ids"`
	TargetAccountNames pq.StringArray `db:"target_account_names"`
	Requester          sql.NullString `db:"requester"`
	Description        sql.NullString `db:"description"`
}
