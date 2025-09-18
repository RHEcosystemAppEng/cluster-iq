package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// CreateExpense represents the data transfer object for creating a new expense.
type ExpenseDTORequest struct {
	InstanceID string    `json:"instanceID"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}

// TODO: comments
func (e ExpenseDTORequest) ToInventoryExpense() *inventory.Expense {
	return inventory.NewExpense(
		e.InstanceID,
		e.Amount,
		e.Date,
	)
}

// TODO comments
type ExpenseDTORequestList struct {
	Expenses []ExpenseDTORequest `json:"expense"`
}

// TODO comments
func NewExpenseDTORequestList(expenses []ExpenseDTOResponse) *ExpenseDTOResponseList {
	Request := ExpenseDTOResponseList{Expenses: expenses}

	// Count only set list length > 0
	if count := len(expenses); count > 0 {
		Request.Count = count
	}

	return &Request
}

// Expense represents the data transfer object for an expense.
type ExpenseDTOResponse struct {
	InstanceID string    `json:"instanceID"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}

// TODO comments
type ExpenseDTOResponseList struct {
	Count    int                  `json:"count,omitempty"`
	Expenses []ExpenseDTOResponse `json:"expense"`
}

// TODO comments
func NewExpenseDTOResponseList(expenses []ExpenseDTOResponse) *ExpenseDTOResponseList {
	response := ExpenseDTOResponseList{Expenses: expenses}

	// Count only set list length > 0
	if count := len(expenses); count > 0 {
		response.Count = count
	}

	return &response
}
