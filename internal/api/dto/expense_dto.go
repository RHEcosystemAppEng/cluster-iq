package dto

import "time"

// Expense represents the data transfer object for an expense.
type Expense struct {
	InstanceID string    `json:"instance_id"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}
