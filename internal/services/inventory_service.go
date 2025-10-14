package services

import (
	"context"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// InventoryService defines the interface for inventory-related business logic.
type InventoryService interface {
	Refresh(ctx context.Context) error
}

var _ InventoryService = (*inventoryServiceImpl)(nil)

type inventoryServiceImpl struct {
	repo repositories.InventoryRepository
	// other dependencies like other services or clients
}

// NewInventoryService creates a new instance of InventoryService.
func NewInventoryService(repo repositories.InventoryRepository) InventoryService {
	return &inventoryServiceImpl{
		repo: repo,
	}
}

// Refresh updates the views and inventory results
func (s *inventoryServiceImpl) Refresh(ctx context.Context) error {
	return s.repo.Refresh(ctx)
}
