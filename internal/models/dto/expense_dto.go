package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// CreateExpense represents the data transfer object for creating a new expense.
type ExpenseDTORequest struct {
	InstanceID string    `json:"instanceId"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
} // @name ExpenseRequest

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

func ToExpenseDTORequest(expense inventory.Expense) *ExpenseDTORequest {
	return &ExpenseDTORequest{
		InstanceID: expense.InstanceID,
		Amount:     expense.Amount,
		Date:       expense.Date,
	}
}

func ToExpenseDTORequestList(expenses []inventory.Expense) *[]ExpenseDTORequest {
	expenseList := make([]ExpenseDTORequest, len(expenses))
	for i, expense := range expenses {
		expenseList[i] = *ToExpenseDTORequest(expense)
	}

	return &expenseList
}

// ExpenseDTOResponse represents the data transfer object for an expense response,
// containing expense details for a specific instance.
type ExpenseDTOResponse struct {
	InstanceID string    `json:"instanceId"`
	Amount     float64   `json:"amount"`
	Date       time.Time `json:"date"`
} // @name ExpenseResponse
