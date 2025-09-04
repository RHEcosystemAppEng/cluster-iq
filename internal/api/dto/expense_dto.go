package dto

import "time"

// Expense represents the data transfer object for an expense.
type Expense struct {
	InstanceID string    `json:"instanceId"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}

// CreateExpense represents the data transfer object for creating a new expense.
type CreateExpense struct {
	InstanceID string    `json:"instanceId" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
}
