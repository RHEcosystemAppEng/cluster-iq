package services

import (
	"context"
	"fmt"

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
	Update(ctx context.Context, action actions.Action) error
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
	action, err := s.repo.GetByID(ctx, actionID)
	if err != nil {
		return action, fmt.Errorf("get action %s: %w", actionID, err)
	}
	return action, nil
}

// Create creates new scheduled actions.
func (s *actionServiceImpl) Create(ctx context.Context, newActions []actions.Action) error {
	if err := s.repo.Create(ctx, newActions); err != nil {
		return fmt.Errorf("create actions: %w", err)
	}
	return nil
}

// Enable enables a scheduled action.
func (s *actionServiceImpl) Enable(ctx context.Context, actionID string) error {
	if err := s.repo.Enable(ctx, actionID); err != nil {
		return fmt.Errorf("enable action %s: %w", actionID, err)
	}
	return nil
}

// Disable disables a scheduled action.
func (s *actionServiceImpl) Disable(ctx context.Context, actionID string) error {
	if err := s.repo.Disable(ctx, actionID); err != nil {
		return fmt.Errorf("disable action %s: %w", actionID, err)
	}
	return nil
}

// Delete removes a scheduled action by its ID.
func (s *actionServiceImpl) Delete(ctx context.Context, actionID string) error {
	if err := s.repo.Delete(ctx, actionID); err != nil {
		return fmt.Errorf("delete action %s: %w", actionID, err)
	}
	return nil
}

// Update updates an action.
func (s *actionServiceImpl) Update(ctx context.Context, action actions.Action) error {
	if err := s.repo.Update(ctx, action); err != nil {
		return fmt.Errorf("update action %s: %w", action.GetID(), err)
	}
	return nil
}
