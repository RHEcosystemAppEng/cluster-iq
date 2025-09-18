package db

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

type ClusterEventDBResponse struct {
	ID             int64     `db:"id"`
	EventTimestamp time.Time `db:"event_timestamp"`
	TriggeredBy    string    `db:"triggered_by"`
	ActionName     string    `db:"action_name"`
	ResourceID     string    `db:"resource_id"`
	ResourceType   string    `db:"resource_type"`
	Result         string    `db:"result"`
	Description    *string   `db:"description,omitempty"`
	Severity       string    `db:"severity"`
}

func (c ClusterEventDBResponse) ToClusterEventDTOResponse() *dto.ClusterEventDTOResponse {
	return &dto.ClusterEventDTOResponse{
		ID:             c.ID,
		EventTimestamp: c.EventTimestamp,
		TriggeredBy:    c.TriggeredBy,
		ActionName:     c.ActionName,
		ResourceID:     c.ResourceID,
		ResourceType:   c.ResourceType,
		Result:         c.Result,
		Description:    c.Description,
		Severity:       c.Severity,
	}
}

type SystemEventDBResponse struct {
	ClusterEventDBResponse
	AccountID string `db:"account_id"`
	Provider  string `db:"provider"`
}

func (s SystemEventDBResponse) ToSystemEventDTOResponse() *dto.SystemEventDTOResponse {
	return &dto.SystemEventDTOResponse{
		ClusterEventDTOResponse: *s.ClusterEventDBResponse.ToClusterEventDTOResponse(),
		AccountID:               s.AccountID,
		Provider:                s.Provider,
	}
}
