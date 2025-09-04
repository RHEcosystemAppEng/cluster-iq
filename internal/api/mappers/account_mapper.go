package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToAccountDTO converts an inventory.Account model to a dto.Account.
func ToAccountDTO(model *inventory.Account) dto.Account {
	if model == nil {
		return dto.Account{}
	}
	return dto.Account{
		ID:                    model.ID,
		Name:                  model.Name,
		Provider:              string(model.Provider),
		ClusterCount:          model.ClusterCount,
		LastScanTimestamp:     model.LastScanTimestamp,
		TotalCost:             model.TotalCost,
		Last15DaysCost:        model.Last15DaysCost,
		LastMonthCost:         model.LastMonthCost,
		CurrentMonthSoFarCost: model.CurrentMonthSoFarCost,
	}
}

// ToAccountDTOList converts a slice of inventory.Account models to a slice of dto.Account.
func ToAccountDTOList(models []inventory.Account) []dto.Account {
	dtos := make([]dto.Account, len(models))
	for i, model := range models {
		dtos[i] = ToAccountDTO(&model)
	}
	return dtos
}

// ToAccountModels converts a slice of dto.NewAccount to a slice of inventory.Account models.
func ToAccountModels(dtos []dto.NewAccount) []inventory.Account {
	models := make([]inventory.Account, len(dtos))
	for i, newAccountDTO := range dtos {
		provider := inventory.GetCloudProvider(newAccountDTO.Provider)
		model := inventory.NewAccount(newAccountDTO.ID, newAccountDTO.Name, provider, "", "")
		models[i] = *model
	}
	return models
}
