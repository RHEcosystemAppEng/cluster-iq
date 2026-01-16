package services

import (
	"context"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// EventService defines the interface for event-related business logic.
type EventService interface {
	ListSystemEvents(ctx context.Context, opts models.ListOptions) ([]db.SystemEventDBResponse, int, error)
	ListClusterEvents(ctx context.Context, opts models.ListOptions) ([]db.ClusterEventDBResponse, int, error)
	Create(ctx context.Context, event events.Event) (int64, error)
	Update(ctx context.Context, eventID int64, result string) error
}

var _ EventService = (*eventServiceImpl)(nil)

type eventServiceImpl struct {
	repo repositories.EventRepository
}

// NewEventService creates a new instance of EventService.
func NewEventService(repo repositories.EventRepository) EventService {
	return &eventServiceImpl{
		repo: repo,
	}
}

// ListSystemEvents retrieves a paginated list of system-level audit events.
func (s *eventServiceImpl) ListSystemEvents(ctx context.Context, opts models.ListOptions) ([]db.SystemEventDBResponse, int, error) {
	return s.repo.ListSystemEvents(ctx, opts)
}

// ListClusterEvents retrieves a paginated list of cluster-specific audit events.
func (s *eventServiceImpl) ListClusterEvents(ctx context.Context, opts models.ListOptions) ([]db.ClusterEventDBResponse, int, error) {
	return s.repo.ListClusterEvents(ctx, opts)
}

// Add creates a new audit event.
func (s *eventServiceImpl) Create(ctx context.Context, event events.Event) (int64, error) {
	return s.repo.CreateEvent(ctx, event)
}

// UpdateStatus updates the status of an existing audit event.
func (s *eventServiceImpl) Update(ctx context.Context, eventID int64, result string) error {
	return s.repo.UpdateEventStatus(ctx, eventID, result)
}
