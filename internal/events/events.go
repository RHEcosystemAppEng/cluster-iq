package events

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"go.uber.org/zap"
)

func NewEventService(sqlClient SQLEventClient, logger *zap.Logger) *EventService {
	return &EventService{
		sqlClient: sqlClient,
		logger:    logger,
	}
}

func (e *EventService) LogEvent(opts EventOptions) (int64, error) {
	log := logger.NewLogger()
	event := models.AuditLog{
		TriggeredBy:    opts.TriggeredBy,
		ActionName:     opts.Action,
		ResourceID:     opts.ResourceID,
		ResourceType:   opts.ResourceType,
		Result:         opts.Result,
		Reason:         opts.Reason,
		Severity:       opts.Severity,
		EventTimestamp: time.Now(),
	}
	eventID, err := e.sqlClient.AddEvent(event)
	if err != nil {
		log.Error("Failed to log event", zap.Error(err))
		return 0, err
	}
	return eventID, nil
}

func (e *EventService) UpdateEventStatus(eventID int64, result string) error {
	log := logger.NewLogger()
	err := e.sqlClient.UpdateEventStatus(eventID, result)
	if err != nil {
		log.Error("Failed to update event status", zap.Int64("event_id", eventID), zap.Error(err))
		return err
	}
	return nil
}

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

func (t *EventTracker) Success() {
	if err := t.service.UpdateEventStatus(t.eventID, ResultSuccess); err != nil {
		t.logger.Error("Failed to update event status", zap.Error(err))
	}
}

func (t *EventTracker) Failed() {
	if err := t.service.UpdateEventStatus(t.eventID, ResultFailed); err != nil {
		t.logger.Error("Failed to update event status", zap.Error(err))
	}
}

func ToAuditEvents(logs []models.AuditLog) []AuditEvent {
	events := make([]AuditEvent, len(logs))
	for i, log := range logs {
		events[i] = AuditEvent{
			ID:             log.ID,
			ActionName:     log.ActionName,
			EventTimestamp: log.EventTimestamp,
			Reason:         log.Reason,
			ResourceID:     log.ResourceID,
			ResourceType:   log.ResourceType,
			Result:         log.Result,
			Severity:       log.Severity,
			TriggeredBy:    log.TriggeredBy,
		}
	}
	return events
}
