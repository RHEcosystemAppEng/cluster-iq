package db

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// TODO comments
type TagDBResponse struct {
	Key        string `db:"key"`
	Value      string `db:"value"`
	InstanceID string `db:"instance_id"`
}

// TODO comments
func (t TagDBResponse) ToTagDTOResponse() *dto.TagDTOResponse {
	return &dto.TagDTOResponse{
		Key:        t.Key,
		Value:      t.Value,
		InstanceID: t.InstanceID,
	}
}
