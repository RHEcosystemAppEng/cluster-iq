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

func ToInventoryExpenseList(dtos []ExpenseDTORequest) *[]inventory.Expense {
	expenses := make([]inventory.Expense, len(dtos))
	for i, dto := range dtos {
		expenses[i] = *dto.ToInventoryExpense()
	}

	return &expenses
}

// Expense represents the data transfer object for an expense.
type ExpenseDTOResponse struct {
	InstanceID string    `json:"instanceID"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
}
