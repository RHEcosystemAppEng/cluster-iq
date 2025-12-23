package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

type ClusterEventDBResponse struct {
	ID             int64                `db:"id"`
	EventTimestamp time.Time            `db:"event_timestamp"`
	TriggeredBy    string               `db:"triggered_by"`
	Action         string               `db:"action"`
	ResourceID     string               `db:"resource_id"`
	ResourceType   string               `db:"resource_type"`
	Result         actions.ActionStatus `db:"result"`
	Description    *string              `db:"description,omitempty"`
	Severity       string               `db:"severity"`
}

func (c ClusterEventDBResponse) ToClusterEventDTOResponse() *dto.ClusterEventDTOResponse {
	return &dto.ClusterEventDTOResponse{
		ID:             c.ID,
		EventTimestamp: c.EventTimestamp,
		TriggeredBy:    c.TriggeredBy,
		Action:         c.Action,
		ResourceID:     c.ResourceID,
		ResourceType:   c.ResourceType,
		Result:         c.Result,
		Description:    c.Description,
		Severity:       c.Severity,
	}
}

func ToClusterEventDTOResponseList(models []ClusterEventDBResponse) []dto.ClusterEventDTOResponse {
	dtos := make([]dto.ClusterEventDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = *model.ToClusterEventDTOResponse()
	}

	return dtos
}

type SystemEventDBResponse struct {
	ClusterEventDBResponse
	AccountID string `db:"account_id"`
	Provider  string `db:"provider"`
}

func (s SystemEventDBResponse) ToSystemEventDTOResponse() *dto.SystemEventDTOResponse {
	return &dto.SystemEventDTOResponse{
		ClusterEventDTOResponse: *s.ToClusterEventDTOResponse(),
		AccountID:               s.AccountID,
		Provider:                s.Provider,
	}
}

func ToSystemEventDTOResponseList(models []SystemEventDBResponse) []dto.SystemEventDTOResponse {
	dtos := make([]dto.SystemEventDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = *model.ToSystemEventDTOResponse()
	}

	return dtos
}
