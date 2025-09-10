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

// ToTagModel converts a dto.Tag to an inventory.Tag model.
func ToTagModel(dto dto.Tag) inventory.Tag {
	return inventory.Tag{
		Key:   dto.Key,
		Value: dto.Value,
	}
}

// ToTagsDTOList converts a slice of inventory.Tag models to a slice of dto.Tag.
func ToTagsDTOList(models []inventory.Tag) []dto.Tag {
	if models == nil {
		return []dto.Tag{}
	}
	dtos := make([]dto.Tag, len(models))
	for i, model := range models {
		dtos[i] = ToTagDTO(model)
	}
	return dtos
}

// ToTagsModelList converts a slice of dto.Tag to a slice of inventory.Tag models.
func ToTagsModelList(dtos []dto.Tag) []inventory.Tag {
	if dtos == nil {
		return []inventory.Tag{}
	}
	models := make([]inventory.Tag, len(dtos))
	for i, d := range dtos {
		models[i] = ToTagModel(d)
	}
	return models
}
