package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// ToClusterModel converts a dto.Cluster to an inventory.Cluster model.
func ToClusterModel(object dto.ClusterDTORequest) inventory.Cluster {
	// Note: Instances and Tags are not mapped back from DTO to model
	// as they are typically read-only details in this direction.
	return inventory.Cluster{
		ClusterID:   object.ClusterID,
		ClusterName: object.ClusterName,
		InfraID:     object.InfraID,
		Provider:    object.Provider,
		Status:      object.Status,
		Region:      object.Region,
		AccountID:   object.AccountID,
		ConsoleLink: object.ConsoleLink,
		LastScanTS:  object.LastScanTS,
		CreatedAt:   object.CreatedAt,
		Age:         object.Age,
		Owner:       object.Owner,
	}
}

// ToClusterModelList converts a slice of dto.Cluster to a slice of inventory.Cluster models.
func ToClusterModelList(dtos []dto.ClusterDTORequest) []inventory.Cluster {
	models := make([]inventory.Cluster, len(dtos))
	for i, d := range dtos {
		models[i] = ToClusterModel(d)
	}
	return models
}

// ToClusterDTO converts an inventory.Cluster model to a dto.Cluster.
func ToClusterDTOResponse(model db.ClusterDBResponse) dto.ClusterDTOResponse {
	return dto.ClusterDTOResponse{
		ClusterID:             model.ClusterID,
		ClusterName:           model.ClusterName,
		InfraID:               model.InfraID,
		Provider:              model.Provider,
		Status:                model.Status,
		Region:                model.Region,
		AccountID:             model.AccountID,
		ConsoleLink:           model.ConsoleLink,
		LastScanTS:            model.LastScanTS,
		CreatedAt:             model.CreatedAt,
		Age:                   model.Age,
		Owner:                 model.Owner,
		TotalCost:             model.TotalCost,
		Last15DaysCost:        model.Last15DaysCost,
		LastMonthCost:         model.LastMonthCost,
		CurrentMonthSoFarCost: model.CurrentMonthSoFarCost,
		// Tags are not a direct field in the model, they are part of instances
	}
}

// ToClusterDTOList converts a slice of inventory.Cluster models to a slice of dto.Cluster.
func ToClusterDTOResponseList(models []db.ClusterDBResponse) []dto.ClusterDTOResponse {
	dtos := make([]dto.ClusterDTOResponse, len(models))
	for i := range models {
		dtos[i] = ToClusterDTOResponse(models[i])
	}
	return dtos
}
