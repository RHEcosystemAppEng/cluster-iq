// Package events provides functionality for audit logging and event tracking.
package events

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/audit"
)

// ToAuditEvents converts AuditLogs to AuditEvents
func ToAuditEvents(logs []audit.AuditLog) []AuditEvent {
	events := make([]AuditEvent, len(logs))
	for i, log := range logs {
		events[i] = convertToAuditEvent(log)
	}
	return events
}

// ToSystemAuditEvents converts SystemAuditLogs to SystemAuditEvents
func ToSystemAuditEvents(logs []audit.SystemAuditLogs) []SystemAuditEvent {
	events := make([]SystemAuditEvent, len(logs))
	for i, log := range logs {
		events[i] = SystemAuditEvent{
			AuditEvent: convertToAuditEvent(log.AuditLog),
			AccountID:  log.AccountID,
			Provider:   log.Provider,
		}
	}
	return events
}

// convertToAuditEvent converts single AuditLog to AuditEvent
func convertToAuditEvent(log audit.AuditLog) AuditEvent {
	return AuditEvent{
		ID:             log.ID,
		ActionName:     log.ActionName,
		EventTimestamp: log.EventTimestamp,
		Description:    log.Description,
		ResourceID:     log.ResourceID,
		ResourceType:   log.ResourceType,
		Result:         log.Result,
		Severity:       log.Severity,
		TriggeredBy:    log.TriggeredBy,
	}
}
