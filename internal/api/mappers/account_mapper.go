package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

func ToAccountModel(object dto.AccountDTORequest) inventory.Account {
	return inventory.Account{
		AccountID:   object.AccountID,
		AccountName: object.AccountName,
		Provider:    object.Provider,
		LastScanTS:  object.LastScanTS,
		CreatedAt:   object.CreatedAt,
		Clusters:    make(map[string]*inventory.Cluster, 0),
	}
}

// ToAccountModelList converts a slice of dto.NewAccount to a slice of inventory.Account models.
func ToAccountModelList(dtos []dto.AccountDTORequest) []inventory.Account {
	models := make([]inventory.Account, len(dtos))
	for i, d := range dtos {
		models[i] = ToAccountModel(d)
	}
	return models
}

// ToAccountDTO converts an inventory.Account model to a dto.Account.
func ToAccountDTOResponse(model db.AccountDBResponse) dto.AccountDTOResponse {
	return dto.AccountDTOResponse{
		AccountID:             model.AccountID,
		AccountName:           model.AccountName,
		Provider:              model.Provider,
		ClusterCount:          model.ClusterCount,
		LastScanTS:            model.LastScanTS,
		CreatedAt:             model.CreatedAt,
		TotalCost:             model.TotalCost,
		Last15DaysCost:        model.Last15DaysCost,
		LastMonthCost:         model.LastMonthCost,
		CurrentMonthSoFarCost: model.CurrentMonthSoFarCost,
	}
}

// ToAccountDTOList converts a slice of inventory.Account models to a slice of dto.Account.
func ToAccountDTOResponseList(models []db.AccountDBResponse) []dto.AccountDTOResponse {
	dtos := make([]dto.AccountDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = ToAccountDTOResponse(model)
	}
	return dtos
}
