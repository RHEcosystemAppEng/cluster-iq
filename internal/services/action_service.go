package services

import (
	"context"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// ActionService defines the interface for action-related business logic.
type ActionService interface {
	List(ctx context.Context, options models.ListOptions) ([]db.ActionDBResponse, int, error)
	Get(ctx context.Context, actionID string) (db.ActionDBResponse, error)
	Create(ctx context.Context, newActions []actions.Action) error
	Enable(ctx context.Context, actionID string) error
	Disable(ctx context.Context, actionID string) error
	Delete(ctx context.Context, actionID string) error
}

var _ ActionService = (*actionServiceImpl)(nil)

type actionServiceImpl struct {
	repo repositories.ActionRepository
}

// NewActionService creates a new instance of ActionService.
func NewActionService(repo repositories.ActionRepository) ActionService {
	return &actionServiceImpl{
		repo: repo,
	}
}

// List retrieves a paginated list of scheduled actions.
func (s *actionServiceImpl) List(ctx context.Context, options models.ListOptions) ([]db.ActionDBResponse, int, error) {
	return s.repo.List(ctx, options)
}

// Get retrieves a single scheduled action by its ID.
func (s *actionServiceImpl) Get(ctx context.Context, actionID string) (db.ActionDBResponse, error) {
	return s.repo.GetByID(ctx, actionID)
}

// Create creates new scheduled actions.
func (s *actionServiceImpl) Create(ctx context.Context, newActions []actions.Action) error {
	return s.repo.Create(ctx, newActions)
}

// Enable enables a scheduled action.
func (s *actionServiceImpl) Enable(ctx context.Context, actionID string) error {
	return s.repo.Enable(ctx, actionID)
}

// Disable disables a scheduled action.
func (s *actionServiceImpl) Disable(ctx context.Context, actionID string) error {
	return s.repo.Disable(ctx, actionID)
}

// Delete removes a scheduled action by its ID.
func (s *actionServiceImpl) Delete(ctx context.Context, actionID string) error {
	return s.repo.Delete(ctx, actionID)
}
