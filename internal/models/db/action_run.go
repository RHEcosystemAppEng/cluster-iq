package db

import "database/sql"

// ActionRunDBResponse represents the database schema for action execution history.
type ActionRunDBResponse struct {
	ID         string         `db:"id"`
	ScheduleID string         `db:"schedule_id"`
	StartedAt  sql.NullTime   `db:"started_at"`
	FinishedAt sql.NullTime   `db:"finished_at"`
	Status     string         `db:"status"`
	ErrorMsg   sql.NullString `db:"error_msg"`
}
