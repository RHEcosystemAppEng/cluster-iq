package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TODO comments
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

// TODO comments
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
	Count    int                 `json:"count,omitempty"`
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
