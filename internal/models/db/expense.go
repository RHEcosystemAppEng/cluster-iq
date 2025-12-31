package db

import (
	"time"
)

// ExpenseDBResponse represents the database schema for expense details,
// linking each field to a corresponding column in the database.
type ExpenseDBResponse struct {
	InstanceID string    `db:"instance_id"`
	Amount     float64   `db:"amount"`
	Date       time.Time `db:"date"`
}
