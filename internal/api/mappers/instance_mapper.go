package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToInstanceDTO converts an inventory.Instance model to a dto.Instance.
func ToInstanceDTO(model inventory.Instance) dto.Instance {
	return dto.Instance{
		ID:                model.ID,
		Name:              model.Name,
		InstanceType:      model.InstanceType,
		ClusterID:         model.ClusterID,
		Status:            string(model.Status),
		CreationTimestamp: model.CreationTimestamp,
		Tags:              ToTagDTOs(model.Tags),
	}
}

// ToInstanceDTOs converts a slice of inventory.Instance models to a slice of dto.Instance.
func ToInstanceDTOs(models []inventory.Instance) []dto.Instance {
	dtos := make([]dto.Instance, len(models))
	for i, model := range models {
		dtos[i] = ToInstanceDTO(model)
	}
	return dtos
}
