package eventservice

import (
	"context"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/audit"
	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"go.uber.org/zap"
)

// Event states
const (
	ResultSuccess = "Success"
	ResultFailed  = "Failed"
	ResultPending = "Pending"
)

// Event severity levels
const (
	SeverityInfo    = "Info"
	SeverityError   = "Error"
	SeverityWarning = "Warning"
)

type EventOptions struct {
	Action       actions.ActionOperation
	Description  *string
	ResourceID   string
	ResourceType string
	Result       string
	Severity     string
	TriggeredBy  string
}

// EventService to write events from every clusteriq component
type EventService struct {
	repo   repositories.EventRepository
	logger *zap.Logger
}

// NewEventService creates a new EventService instance.
func NewEventService(dbClient *dbclient.DBClient, logger *zap.Logger) *EventService {
	return &EventService{
		repo:   repositories.NewEventRepository(dbClient),
		logger: logger,
	}
}

// LogEvent creates a new audit log entry and returns its ID.
func (e *EventService) LogEvent(opts EventOptions) (int64, error) {
	event := audit.AuditLog{
		TriggeredBy:    opts.TriggeredBy,
		ActionName:     opts.Action,
		ResourceID:     opts.ResourceID,
		ResourceType:   opts.ResourceType,
		Result:         opts.Result,
		Description:    opts.Description,
		Severity:       opts.Severity,
		EventTimestamp: time.Now().UTC(),
	}
	// TODO Fix replace TODO context by request's context
	eventID, err := e.repo.AddEvent(context.TODO(), event)
	if err != nil {
		e.logger.Error("Failed to log event", zap.Error(err))
		return 0, err
	}
	return eventID, nil
}

// UpdateEventStatus updates the result status of an existing event.
func (e *EventService) UpdateEventStatus(eventID int64, result string) error {
	// TODO Fix replace TODO context by request's context
	err := e.repo.UpdateEventStatus(context.TODO(), eventID, result)
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

type EventTracker struct {
	eventID int64
	service *EventService
	logger  *zap.Logger
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
