package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

// ToClusterEventDTO converts an audit.AuditLog model to a dto.ClusterEvent.
func ToClusterEventDTOResponse(model db.ClusterEventDBResponse) dto.ClusterEventDTOResponse {
	return dto.ClusterEventDTOResponse{
		ID:             model.ID,
		ActionName:     string(model.ActionName),
		EventTimestamp: model.EventTimestamp,
		Description:    model.Description,
		ResourceID:     model.ResourceID,
		ResourceType:   model.ResourceType,
		Result:         model.Result,
		Severity:       model.Severity,
		TriggeredBy:    model.TriggeredBy,
	}
}

// ToClusterEventDTOList converts a slice of audit.AuditLog models to a slice of dto.ClusterEvent.
func ToClusterEventDTOResponseList(models []db.ClusterEventDBResponse) []dto.ClusterEventDTOResponse {
	dtos := make([]dto.ClusterEventDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = ToClusterEventDTOResponse(model)
	}
	return dtos
}

// ToSystemEventDTO converts an audit.SystemAuditLogs model to a dto.SystemEvent.
func ToSystemEventDTOResponse(model db.SystemEventDBResponse) dto.SystemEventDTOResponse {
	return dto.SystemEventDTOResponse{
		ClusterEventDTOResponse: ToClusterEventDTOResponse(model.ClusterEventDBResponse),
		AccountID:               model.AccountID,
		Provider:                model.Provider,
	}
}

// ToSystemEventDTOList converts a slice of audit.SystemAuditLogs models to a slice of dto.SystemEvent.
func ToSystemEventDTOList(models []db.SystemEventDBResponse) []dto.SystemEventDTOResponse {
	dtos := make([]dto.SystemEventDTOResponse, len(models))
	for i, model := range models {
		dtos[i] = ToSystemEventDTOResponse(model)
	}
	return dtos
}
