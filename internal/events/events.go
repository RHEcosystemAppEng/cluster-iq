// Package events provides functionality for audit logging and event tracking.
package events

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"go.uber.org/zap"
)

// NewEventService creates a new EventService instance.
func NewEventService(sqlClient SQLEventClient, logger *zap.Logger) *EventService {
	return &EventService{
		sqlClient: sqlClient,
		logger:    logger,
	}
}

// LogEvent creates a new audit log entry and returns its ID.
func (e *EventService) LogEvent(opts EventOptions) (int64, error) {
	event := models.AuditLog{
		TriggeredBy:    opts.TriggeredBy,
		ActionName:     opts.Action,
		ResourceID:     opts.ResourceID,
		ResourceType:   opts.ResourceType,
		Result:         opts.Result,
		Description:    opts.Description,
		Severity:       opts.Severity,
		EventTimestamp: time.Now().UTC(),
	}
	eventID, err := e.sqlClient.AddEvent(event)
	if err != nil {
		e.logger.Error("Failed to log event", zap.Error(err))
		return 0, err
	}
	return eventID, nil
}

// UpdateEventStatus updates the result status of an existing event.
func (e *EventService) UpdateEventStatus(eventID int64, result string) error {
	err := e.sqlClient.UpdateEventStatus(eventID, result)
	if err != nil {
		e.logger.Error("Failed to update event status", zap.Int64("event_id", eventID), zap.Error(err))
		return err
	}
	return nil
}

// StartTracking begins tracking a new event and returns an EventTracker.
func (e *EventService) StartTracking(opts *EventOptions) *EventTracker {
	eventID, err := e.LogEvent(*opts)
	if err != nil {
		e.logger.Error("Failed to log initial event", zap.Error(err))
		return nil
	}

	return &EventTracker{
		eventID: eventID,
		service: e,
		logger:  e.logger,
	}
}

// Failed marks the tracked event status as failed.
func (t *EventTracker) Success() {
	if err := t.service.UpdateEventStatus(t.eventID, ResultSuccess); err != nil {
		t.logger.Error("Failed to update event status", zap.Error(err))
	}
}

// Failed marks the tracked event as failed.
func (t *EventTracker) Failed() {
	if err := t.service.UpdateEventStatus(t.eventID, ResultFailed); err != nil {
		t.logger.Error("Failed to update event status", zap.Error(err))
	}
}

// ToAuditEvents converts AuditLogs to AuditEvents
func ToAuditEvents(logs []models.AuditLog) []AuditEvent {
	events := make([]AuditEvent, len(logs))
	for i, log := range logs {
		events[i] = convertToAuditEvent(log)
	}
	return events
}

// ToSystemAuditEvents converts SystemAuditLogs to SystemAuditEvents
func ToSystemAuditEvents(logs []models.SystemAuditLogs) []SystemAuditEvent {
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
func convertToAuditEvent(log models.AuditLog) AuditEvent {
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
