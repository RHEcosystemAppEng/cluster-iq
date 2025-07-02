package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToTagDTO converts an inventory.Tag model to a dto.Tag.
func ToTagDTO(model inventory.Tag) dto.Tag {
	return dto.Tag{
		Key:   model.Key,
		Value: model.Value,
	}
}

// ToTagDTOs converts a slice of inventory.Tag models to a slice of dto.Tag.
func ToTagDTOs(models []inventory.Tag) []dto.Tag {
	if models == nil {
		return []dto.Tag{}
	}
	dtos := make([]dto.Tag, len(models))
	for i, model := range models {
		dtos[i] = ToTagDTO(model)
	}
	return dtos
}
