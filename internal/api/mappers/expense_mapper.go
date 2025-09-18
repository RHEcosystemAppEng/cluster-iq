package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// ToExpenseModel converts a dto.CreateExpense DTO to an inventory.Expense model.
func ToExpenseModel(object dto.ExpenseDTORequest) inventory.Expense {
	return inventory.Expense{
		InstanceID: object.InstanceID,
		Amount:     object.Amount,
		Date:       object.Date,
	}
}

func ToExpenseModelList(dtos []dto.ExpenseDTORequest) []inventory.Expense {
	models := make([]inventory.Expense, len(dtos))
	for i, d := range dtos {
		models[i] = ToExpenseModel(d)
	}
	return models
}

// ToExpenseDTO converts an inventory.Expense model to a dto.Expense DTO.
func ToExpenseDTOResponse(model db.ExpenseDBResponse) dto.ExpenseDTOResponse {
	return dto.ExpenseDTOResponse{
		InstanceID: model.InstanceID,
		Amount:     model.Amount,
		Date:       model.Date,
	}
}

// ToExpenseDTOList converts a slice of inventory.Expense models to a slice of dto.Expense DTOs.
func ToExpenseDTOList(models []db.ExpenseDBResponse) []dto.ExpenseDTOResponse {
	dtos := make([]dto.ExpenseDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = ToExpenseDTOResponse(model)
	}
	return dtos
}
