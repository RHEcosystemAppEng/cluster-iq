package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToInstanceDTO converts an inventory.Instance model to a dto.Instance.
func ToInstanceDTO(model *inventory.Instance) dto.Instance {
	if model == nil {
		return dto.Instance{}
	}
	return dto.Instance{
		ID:                model.ID,
		Name:              model.Name,
		Provider:          string(model.Provider),
		InstanceType:      model.InstanceType,
		ClusterID:         model.ClusterID,
		Status:            string(model.Status),
		CreationTimestamp: model.CreationTimestamp,
		Tags:              ToTagsDTOList(model.Tags),
	}
}

// ToInstanceModel converts a dto.Instance to an inventory.Instance model.
func ToInstanceModel(dto dto.Instance) inventory.Instance {
	return inventory.Instance{
		ID:                dto.ID,
		Provider:          inventory.CloudProvider(dto.Provider),
		InstanceType:      dto.InstanceType,
		AvailabilityZone:  dto.AvailabilityZone,
		ClusterID:         dto.ClusterID,
		Status:            inventory.InstanceStatus(dto.Status),
		Tags:              ToTagsModelList(dto.Tags),
		CreationTimestamp: dto.CreationTimestamp,
	}
}

// ToInstanceModelList converts a slice of dto.Instance to a slice of inventory.Instance models.
func ToInstanceModelList(dtos []dto.Instance) []inventory.Instance {
	models := make([]inventory.Instance, len(dtos))
	for i, d := range dtos {
		models[i] = ToInstanceModel(d)
	}
	return models
}

// ToInstanceDTOList converts a slice of inventory.Instance models to a slice of dto.Instance.
func ToInstanceDTOList(models []inventory.Instance) []dto.Instance {
	dtos := make([]dto.Instance, len(models))
	for i, model := range models {
		dtos[i] = ToInstanceDTO(&model)
	}
	return dtos
}
