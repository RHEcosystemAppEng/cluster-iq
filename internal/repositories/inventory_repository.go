package repositories

import (
	"context"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
)

var _ InventoryRepository = (*inventoryRepositoryImpl)(nil)

// InventoryRepository defines the interface for data access operations for inventory.
type InventoryRepository interface {
	Refresh(ctx context.Context) error
}

type inventoryRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewInventoryRepository(db *dbclient.DBClient) InventoryRepository {
	return &inventoryRepositoryImpl{db: db}
}

func (r *inventoryRepositoryImpl) Refresh(ctx context.Context) error {
	if err := r.db.ExecFunc(RefreshMaterializedViewsQuery); err != nil {
		return err
	}

	return nil
}
