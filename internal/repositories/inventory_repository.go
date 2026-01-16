package repositories

import (
	"context"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
)

const (
	// Run 'refresh_materialized_views()' SQL function
	RefreshMaterializedViewsQuery = "SELECT refresh_materialized_views()"
	// Run 'check_terminated_instances()' SQL function
	UpdateTerminatedInstancesQuery = `SELECT check_terminated_instances()`
	// Run 'check_terminated_clusters()' SQL function
	UpdateTerminatedClustersQuery = `SELECT check_terminated_clusters()`
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
	// Updating 'Terminated' Instances
	if err := r.db.ExecFunc(ctx, UpdateTerminatedClustersQuery); err != nil {
		return err
	}

	// Updating 'Terminated' Clusters
	if err := r.db.ExecFunc(ctx, UpdateTerminatedClustersQuery); err != nil {
		return err
	}

	// Refreshing Materialized Views
	if err := r.db.ExecFunc(ctx, RefreshMaterializedViewsQuery); err != nil {
		return err
	}

	return nil
}
