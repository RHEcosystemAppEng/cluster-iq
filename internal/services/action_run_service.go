package services

import (
	"context"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// ActionRunService defines the interface for action run business logic.
type ActionRunService interface {
	List(ctx context.Context, options models.ListOptions) ([]db.ActionRunDBResponse, int, error)
	Get(ctx context.Context, runID string) (db.ActionRunDBResponse, error)
	Create(ctx context.Context, scheduleID string) (int64, error)
	Update(ctx context.Context, runID string, status string, errorMsg string) error
}

var _ ActionRunService = (*actionRunServiceImpl)(nil)

type actionRunServiceImpl struct {
	repo repositories.ActionRunRepository
}

func NewActionRunService(repo repositories.ActionRunRepository) ActionRunService {
	return &actionRunServiceImpl{repo: repo}
}

func (s *actionRunServiceImpl) List(ctx context.Context, options models.ListOptions) ([]db.ActionRunDBResponse, int, error) {
	return s.repo.List(ctx, options)
}

func (s *actionRunServiceImpl) Get(ctx context.Context, runID string) (db.ActionRunDBResponse, error) {
	run, err := s.repo.GetByID(ctx, runID)
	if err != nil {
		return run, fmt.Errorf("get action run %s: %w", runID, err)
	}
	return run, nil
}

func (s *actionRunServiceImpl) Create(ctx context.Context, scheduleID string) (int64, error) {
	id, err := s.repo.Create(ctx, scheduleID)
	if err != nil {
		return -1, fmt.Errorf("create action run: %w", err)
	}
	return id, nil
}

func (s *actionRunServiceImpl) Update(ctx context.Context, runID string, status string, errorMsg string) error {
	if err := s.repo.Update(ctx, runID, status, errorMsg); err != nil {
		return fmt.Errorf("update action run %s: %w", runID, err)
	}
	return nil
}
