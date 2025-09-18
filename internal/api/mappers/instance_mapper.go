package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// ToInstanceModel converts a dto.InstanceDTOResponse to an inventory.Instance model.
func ToInstanceModel(dto dto.InstanceDTORequest) inventory.Instance {
	return inventory.Instance{
		InstanceID:       dto.InstanceID,
		InstanceName:     dto.InstanceName,
		Provider:         dto.Provider,
		InstanceType:     dto.InstanceType,
		AvailabilityZone: dto.AvailabilityZone,
		ClusterID:        dto.ClusterID,
		Status:           dto.Status,
		Tags:             ToTagsModelList(dto.Tags),
		CreatedAt:        dto.CreatedAt,
	}
}

// ToInstanceModelList converts a slice of dto.InstanceDTOResponse to a slice of inventory.Instance models.
func ToInstanceModelList(dtos []dto.InstanceDTORequest) []inventory.Instance {
	models := make([]inventory.Instance, len(dtos))
	for i, d := range dtos {
		models[i] = ToInstanceModel(d)
	}
	return models
}

// ToInstanceDTO converts an inventory.Instance model to a dto.Instance.
func ToInstanceDTOResponse(model db.InstanceDBResponse) dto.InstanceDTOResponse {
	return dto.InstanceDTOResponse{
		InstanceID:   model.InstanceID,
		InstanceName: model.InstanceName,
		Provider:     model.Provider,
		InstanceType: model.InstanceType,
		ClusterID:    model.ClusterID,
		Status:       model.Status,
		CreatedAt:    model.CreatedAt,
		Tags:         ToTagsDTOResponseList(model.Tags),
	}
}

// ToInstanceDTOList converts a slice of inventory.Instance models to a slice of dto.InstanceDTOResponse.
func ToInstanceDTOResponseList(models []db.InstanceDBResponse) []dto.InstanceDTOResponse {
	dtos := make([]dto.InstanceDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = ToInstanceDTOResponse(model)
	}
	return dtos
}
