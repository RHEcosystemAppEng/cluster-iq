package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/jmoiron/sqlx"
)

var _ ClusterRepository = (*clusterRepositoryImpl)(nil)

type ClusterRepository interface {
	GetClusterByID(ctx context.Context, clusterID string) (*inventory.Cluster, error)
	ListClusters(ctx context.Context, opts ListOptions) ([]inventory.Cluster, int, error)
	GetClustersOverview(ctx context.Context) (inventory.ClustersSummary, error)
	GetClusterAccountName(ctx context.Context, clusterID string) (string, error)
	GetClusterRegion(ctx context.Context, clusterID string) (string, error)
	GetClusterTags(ctx context.Context, clusterID string) ([]inventory.Tag, error)
	GetClustersOnAccount(ctx context.Context, accountName string) ([]inventory.Cluster, error)
	GetInstancesOnCluster(ctx context.Context, clusterID string) ([]inventory.Instance, error)
	DeleteCluster(ctx context.Context, id string) error
	CreateClusters(ctx context.Context, clusters []inventory.Cluster) error
	UpdateCluster(ctx context.Context, cluster inventory.Cluster) error
	UpdateClusterStatusByClusterID(ctx context.Context, status string, clusterID string) error
}

type clusterRepositoryImpl struct {
	db *sqlx.DB
}

func NewClusterRepository(db *sqlx.DB) ClusterRepository {
	return &clusterRepositoryImpl{db: db}
}

// GetClusterByID retrieves a cluster's details by its unique identifier.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice containing a single inventory.Cluster object.
// - An error if the query fails or the cluster ID does not exist.
func (r *clusterRepositoryImpl) GetClusterByID(ctx context.Context, clusterID string) (*inventory.Cluster, error) {
	var cluster inventory.Cluster
	err := r.db.GetContext(ctx, &cluster, SelectClustersByIDQuery, clusterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cluster, nil
}

// ListClusters retrieves all clusters from the database.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (r *clusterRepositoryImpl) ListClusters(ctx context.Context, opts ListOptions) ([]inventory.Cluster, int, error) {
	var clusters []inventory.Cluster
	baseQuery := SelectClustersQuery
	countQuery := "SELECT COUNT(*) FROM clusters"

	whereClauses, namedArgs := buildClusterWhereClauses(opts.Filters)

	total, err := listQueryHelper(ctx, r.db, &clusters, baseQuery, countQuery, opts, whereClauses, namedArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list clusters: %w", err)
	}
	return clusters, total, nil
}

// GetClustersOverview returns a summary of cluster statuses
// It counts the number of clusters that are running, stopped or terminated.
func (r *clusterRepositoryImpl) GetClustersOverview(ctx context.Context) (inventory.ClustersSummary, error) {
	var countsDB inventory.ClustersSummary
	err := r.db.GetContext(ctx, &countsDB, SelectClustersOverview)
	if err != nil {
		return inventory.ClustersSummary{}, err
	}
	return countsDB, nil
}

// GetClusterAccountName retrieves the account name associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A string representing the account name.
// - An error if the query fails or the cluster ID does not exist.
func (r *clusterRepositoryImpl) GetClusterAccountName(ctx context.Context, clusterID string) (string, error) {
	var accountName string
	err := r.db.GetContext(ctx, &accountName, SelectClusterAccountNameQuery, clusterID)
	return accountName, err
}

// GetClusterRegion retrieves the region where a specific cluster is located.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A string representing the region of the cluster.
// - An error if the query fails or the cluster ID does not exist.
func (r *clusterRepositoryImpl) GetClusterRegion(ctx context.Context, clusterID string) (string, error) {
	var region string
	err := r.db.GetContext(ctx, &region, SelectClusterRegionQuery, clusterID)
	return region, err
}

// GetClusterTags retrieves the tags associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice of inventory.Tag objects representing the cluster's tags.
// - An error if the query fails.
func (r *clusterRepositoryImpl) GetClusterTags(ctx context.Context, clusterID string) ([]inventory.Tag, error) {
	var tags []inventory.Tag
	err := r.db.SelectContext(ctx, &tags, SelectClusterTags, clusterID)
	return tags, err
}

// GetClustersOnAccount retrieves all clusters associated with a specific account.
//
// Parameters:
// - accountName: The name of the account whose clusters will be retrieved.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (r *clusterRepositoryImpl) GetClustersOnAccount(ctx context.Context, accountName string) ([]inventory.Cluster, error) {
	var clusters []inventory.Cluster
	err := r.db.SelectContext(ctx, &clusters, SelectClustersOnAccountQuery, accountName)
	return clusters, err
}

// GetInstancesOnCluster retrieves all instances belonging to a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice of inventory.Instance objects representing the instances in the cluster.
// - An error if the query fails.
func (r *clusterRepositoryImpl) GetInstancesOnCluster(ctx context.Context, clusterID string) ([]inventory.Instance, error) {
	var instances []inventory.Instance
	err := r.db.SelectContext(ctx, &instances, SelectInstancesOnClusterQuery, clusterID)
	return instances, err
}

// DeleteCluster deletes a cluster from the database.
//
// Parameters:
// - id: The id of the cluster to delete.
//
// Returns:
// - An error if the database transaction fails.
func (r *clusterRepositoryImpl) DeleteCluster(ctx context.Context, id string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a rollback in case anything fails.
	defer func() { _ = tx.Rollback() }() // Rollback is a no-op if the transaction is already committed.

	_, err = tx.ExecContext(ctx, DeleteClusterQuery, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// CreateClusters inserts a list of clusters into the database in a transaction.
//
// Parameters:
// - clusters: A slice of inventory.Cluster objects to insert.
//
// Returns:
// - An error if the transaction fails or the query encounters an issue.
func (r *clusterRepositoryImpl) CreateClusters(ctx context.Context, clusters []inventory.Cluster) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareNamedContext(ctx, InsertClustersQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare named statement: %w", err)
	}
	defer stmt.Close()

	for _, cluster := range clusters {
		if _, err := stmt.ExecContext(ctx, cluster); err != nil {
			return fmt.Errorf("failed to execute statement for cluster %s: %w", cluster.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateCluster updates an existing cluster's details in the database.
func (r *clusterRepositoryImpl) UpdateCluster(ctx context.Context, cluster inventory.Cluster) error {
	_, err := r.db.NamedExecContext(ctx, UpdateClusterQuery, &cluster)
	if err != nil {
		return fmt.Errorf("failed to update cluster %s: %w", cluster.ID, err)
	}
	return nil
}

// UpdateClusterStatusByClusterID updates the status of a cluster and all its instances in the database.
//
// This function first verifies if the requested status exists in the database. If the status is valid, it updates:
// 1. The status of the cluster identified by the given `clusterID`.
// 2. The status of all instances associated with the cluster.
//
// Parameters:
// - status: The new status to be applied to the cluster and its instances.
// - clusterID: The unique identifier of the cluster whose status will be updated.
//
// Returns:
// - An error if the status is invalid, the update operation fails, or no rows are affected.
func (r *clusterRepositoryImpl) UpdateClusterStatusByClusterID(ctx context.Context, status string, clusterID string) error {
	_, err := r.db.ExecContext(ctx, UpdateStatusClusterByClusterIDQuery, status, clusterID)
	return err
}

func buildClusterWhereClauses(filters map[string]interface{}) ([]string, map[string]interface{}) {
	clauses := make([]string, 0, len(filters))
	args := make(map[string]interface{})

	for key, value := range filters {
		switch key {
		case "provider":
			clauses = append(clauses, "provider = :provider")
			args["provider"] = value
		case "region":
			clauses = append(clauses, "region = :region")
			args["region"] = value
		case "status":
			clauses = append(clauses, "status = :status")
			args["status"] = value
		case "account_name":
			clauses = append(clauses, "account_name = :account_name")
			args["account_name"] = value
		}
	}

	return clauses, args
}
