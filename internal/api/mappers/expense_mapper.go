package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToExpenseDTO converts an inventory.Expense model to a dto.Expense.
func ToExpenseDTO(model inventory.Expense) dto.Expense {
	return dto.Expense{
		InstanceID: model.InstanceID,
		Amount:     model.Amount,
		Date:       model.Date,
	}
}

// ToExpenseDTOs converts a slice of inventory.Expense models to a slice of dto.Expense.
func ToExpenseDTOs(models []inventory.Expense) []dto.Expense {
	dtos := make([]dto.Expense, len(models))
	for i, model := range models {
		dtos[i] = ToExpenseDTO(model)
	}
	return dtos
}
