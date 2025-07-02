package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

func ToClusterDTO(model inventory.Cluster) dto.ClusterDTO {
	return dto.ClusterDTO{
		ID:                    model.ID,
		Name:                  model.Name,
		InfraID:               model.InfraID,
		Provider:              model.Provider,
		Status:                model.Status,
		Region:                model.Region,
		AccountName:           model.AccountName,
		ConsoleLink:           model.ConsoleLink,
		InstanceCount:         model.InstanceCount,
		LastScanTimestamp:     model.LastScanTimestamp,
		CreationTimestamp:     model.CreationTimestamp,
		Age:                   model.Age,
		Owner:                 model.Owner,
		TotalCost:             model.TotalCost,
		Last15DaysCost:        model.Last15DaysCost,
		LastMonthCost:         model.LastMonthCost,
		CurrentMonthSoFarCost: model.CurrentMonthSoFarCost,
		// TODO
		Instances: nil,
	}
}

func ToClusterDTOs(models []inventory.Cluster) []dto.ClusterDTO {
	dtos := make([]dto.ClusterDTO, len(models))
	for i, model := range models {
		dtos[i] = ToClusterDTO(model)
	}
	return dtos
}
