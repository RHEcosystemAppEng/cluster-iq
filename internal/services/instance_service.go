package services

import (
	"context"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// InstanceService defines the interface for instance-related business logic.
type InstanceService interface {
	List(ctx context.Context, opts models.ListOptions) ([]db.InstanceDBResponse, int, error)
	Get(ctx context.Context, instanceID string) (db.InstanceDBResponse, error)
	GetSummary(ctx context.Context) (inventory.InstancesSummary, error)
	Create(ctx context.Context, instances []inventory.Instance) error
}

var _ InstanceService = (*instanceServiceImpl)(nil)

type instanceServiceImpl struct {
	repo repositories.InstanceRepository
	// other dependencies like other services or clients
}

// NewInstanceService creates a new instance of InstanceService.
func NewInstanceService(repo repositories.InstanceRepository) InstanceService {
	return &instanceServiceImpl{
		repo: repo,
	}
}

// List retrieves a paginated list of instances.
func (s *instanceServiceImpl) List(ctx context.Context, opts models.ListOptions) ([]db.InstanceDBResponse, int, error) {
	return s.repo.ListInstances(ctx, opts)
}

// Get retrieves a single instance by its ID.
func (s *instanceServiceImpl) Get(ctx context.Context, instanceID string) (db.InstanceDBResponse, error) {
	instance, err := s.repo.GetInstanceByID(ctx, instanceID)
	if err != nil {
		return instance, fmt.Errorf("get instance %s: %w", instanceID, err)
	}
	return instance, nil
}

// GetSummary retrieves a summary of instance counts.
func (s *instanceServiceImpl) GetSummary(ctx context.Context) (inventory.InstancesSummary, error) {
	summary, err := s.repo.GetInstancesOverview(ctx)
	if err != nil {
		return summary, fmt.Errorf("get instances summary: %w", err)
	}
	return summary, nil
}

// Create creates new instances.
func (s *instanceServiceImpl) Create(ctx context.Context, instances []inventory.Instance) error {
	if err := s.repo.CreateInstances(ctx, instances); err != nil {
		return fmt.Errorf("create instances: %w", err)
	}
	return nil
}
