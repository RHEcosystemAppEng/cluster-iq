package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToExpenseDTO converts an inventory.Expense model to a dto.Expense DTO.
func ToExpenseDTO(model inventory.Expense) dto.Expense {
	return dto.Expense{
		InstanceID: model.InstanceID,
		Amount:     model.Amount,
		Date:       model.Date,
	}
}

// ToExpenseModel converts a dto.CreateExpense DTO to an inventory.Expense model.
func ToExpenseModel(dto dto.CreateExpense) inventory.Expense {
	return inventory.Expense{
		InstanceID: dto.InstanceID,
		Amount:     dto.Amount,
		Date:       dto.Date,
	}
}

// ToExpenseDTOList converts a slice of inventory.Expense models to a slice of dto.Expense DTOs.
func ToExpenseDTOList(models []inventory.Expense) []dto.Expense {
	dtos := make([]dto.Expense, len(models))
	for i, model := range models {
		dtos[i] = ToExpenseDTO(model)
	}
	return dtos
}
