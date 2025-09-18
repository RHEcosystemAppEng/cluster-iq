package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// ToTagModel converts a dto.TagDTOResponse to an inventory.Tag model.
func ToTagModel(dto dto.TagDTORequest) inventory.Tag {
	return inventory.Tag{
		Key:   dto.Key,
		Value: dto.Value,
	}
}

// ToTagsModelList converts a slice of dto.TagDTOResponseList to a slice of inventory.Tag models.
func ToTagsModelList(dtos []dto.TagDTORequest) []inventory.Tag {
	models := make([]inventory.Tag, len(dtos))
	for i, d := range dtos {
		models[i] = ToTagModel(d)
	}
	return models
}

// ToTagDTO converts an inventory.Tag model to a dto.TagDTOResponse.
func ToTagDTOResponse(model db.TagDBResponse) dto.TagDTOResponse {
	return dto.TagDTOResponse{
		Key:   model.Key,
		Value: model.Value,
	}
}

// ToTagsDTOList converts a slice of inventory.Tag models to a slice of dto.TagDTOResponse.
func ToTagsDTOResponseList(models []db.TagDBResponse) []dto.TagDTOResponse {
	dtos := make([]dto.TagDTOResponse, len(models))
	for i := range models {
		dtos[i] = ToTagDTOResponse(models[i])
	}
	return dtos
}
