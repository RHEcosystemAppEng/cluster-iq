package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/audit"
)

// ToClusterEventDTO converts an audit.AuditLog model to a dto.ClusterEvent.
func ToClusterEventDTO(model audit.AuditLog) dto.ClusterEvent {
	return dto.ClusterEvent{
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

// ToClusterEventDTOs converts a slice of audit.AuditLog models to a slice of dto.ClusterEvent.
func ToClusterEventDTOs(models []audit.AuditLog) []dto.ClusterEvent {
	dtos := make([]dto.ClusterEvent, len(models))
	for i, model := range models {
		dtos[i] = ToClusterEventDTO(model)
	}
	return dtos
}

// ToSystemEventDTO converts an audit.SystemAuditLogs model to a dto.SystemEvent.
func ToSystemEventDTO(model audit.SystemAuditLogs) dto.SystemEvent {
	return dto.SystemEvent{
		ClusterEvent: ToClusterEventDTO(model.AuditLog),
		AccountID:    model.AccountID,
		Provider:     model.Provider,
	}
}

// ToSystemEventDTOs converts a slice of audit.SystemAuditLogs models to a slice of dto.SystemEvent.
func ToSystemEventDTOs(models []audit.SystemAuditLogs) []dto.SystemEvent {
	dtos := make([]dto.SystemEvent, len(models))
	for i, model := range models {
		dtos[i] = ToSystemEventDTO(model)
	}
	return dtos
}
