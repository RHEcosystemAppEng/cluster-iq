package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToClusterDTO converts an inventory.Cluster model to a dto.Cluster.
func ToClusterDTO(model inventory.Cluster) dto.Cluster {
	return dto.Cluster{
		ID:                    model.ID,
		Name:                  model.Name,
		InfraID:               model.InfraID,
		Provider:              string(model.Provider),
		Status:                string(model.Status),
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
		Instances:             ToInstanceDTOs(model.Instances),
		// Tags are not a direct field in the model, they are part of instances
	}
}

// ToClusterDTOs converts a slice of inventory.Cluster models to a slice of dto.Cluster.
func ToClusterDTOs(models []inventory.Cluster) []dto.Cluster {
	dtos := make([]dto.Cluster, len(models))
	for i, model := range models {
		dtos[i] = ToClusterDTO(model)
	}
	return dtos
}

// ToClusterModel converts a dto.Cluster to an inventory.Cluster model.
func ToClusterModel(dto dto.Cluster) inventory.Cluster {
	// Note: Instances and Tags are not mapped back from DTO to model
	// as they are typically read-only details in this direction.
	return inventory.Cluster{
		ID:          dto.ID,
		Name:        dto.Name,
		InfraID:     dto.InfraID,
		Provider:    inventory.GetProvider(dto.Provider),
		Status:      inventory.AsInstanceStatus(dto.Status),
		Region:      dto.Region,
		AccountName: dto.AccountName,
		ConsoleLink: dto.ConsoleLink,
		Owner:       dto.Owner,
	}
}
