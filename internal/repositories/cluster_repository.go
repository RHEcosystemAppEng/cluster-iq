package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

const (
	// DB Table for clusters
	ClustersTable = "clusters"
	// View for getting clusters
	SelectClustersFullView = "clusters_full_view"
	// View for getting cluster tags
	SelectClusterTags = "clusters_tags"
	// Materialized view for getting clusters
	SelectClustersFullMView = "m_clusters_full_view"
	// InsertClustersQuery to insert or update new clusters
	InsertClustersQuery = `
		INSERT INTO clusters (
			cluster_id,
			cluster_name,
			infra_id,
			provider,
			status,
			region,
			account_id,
			console_link,
			last_scan_ts,
			created_at,
			age,
			owner
		) VALUES (
			:cluster_id,
			:cluster_name,
			:infra_id,
			:provider,
			:status,
			:region,
			(SELECT id FROM accounts WHERE account_id = :account_id),
			:console_link,
			:last_scan_ts,
			:created_at,
			:age,
			:owner
		) ON CONFLICT ON CONSTRAINT uq_clusters_accountid_clusterid DO UPDATE SET
			status = EXCLUDED.status,
			region = EXCLUDED.region,
			console_link = EXCLUDED.console_link,
			last_scan_ts = EXCLUDED.last_scan_ts,
			created_at = EXCLUDED.created_at,
			age = EXCLUDED.age,
			owner = EXCLUDED.owner
	`
)

var _ ClusterRepository = (*clusterRepositoryImpl)(nil)

type ClusterRepository interface {
	ListClusters(ctx context.Context, opts models.ListOptions) ([]db.ClusterDBResponse, int, error)
	CountClusters(ctx context.Context, opts models.ListOptions) (int, error)
	GetClusterByID(ctx context.Context, clusterID string) (*db.ClusterDBResponse, error)
	GetClusterAccountName(ctx context.Context, clusterID string) (string, error)
	GetClusterRegion(ctx context.Context, clusterID string) (string, error)
	GetClusterTags(ctx context.Context, clusterID string) ([]db.TagDBResponse, error)
	GetClustersOnAccount(ctx context.Context, accountName string) ([]db.ClusterDBResponse, error)
	GetInstancesOnCluster(ctx context.Context, clusterID string) ([]db.InstanceDBResponse, error)
	GetClustersOverview(ctx context.Context) (inventory.ClustersSummary, error)
	CreateClusters(ctx context.Context, clusters []inventory.Cluster) error
	UpdateCluster(ctx context.Context, clusterID string, patch dto.ClusterPatchRequest) error
	UpdateClusterStatusByClusterID(ctx context.Context, status string, clusterID string) error
	DeleteCluster(ctx context.Context, id string) error
}

type clusterRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewClusterRepository(db *dbclient.DBClient) ClusterRepository {
	return &clusterRepositoryImpl{db: db}
}

// ListClusters retrieves all clusters from the database.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (r *clusterRepositoryImpl) ListClusters(ctx context.Context, opts models.ListOptions) ([]db.ClusterDBResponse, int, error) {
	clusters := []db.ClusterDBResponse{}

	if err := r.db.SelectWithContext(ctx, &clusters, SelectClustersFullMView, opts, "cluster_id", "*"); err != nil {
		return clusters, 0, fmt.Errorf("failed to list clusters: %w", err)
	}

	return clusters, len(clusters), nil
}

// CountClusters performs counting operations.
//
// Returns:
// - The number of counted clusters.
// - An error if the query fails.
func (r *clusterRepositoryImpl) CountClusters(ctx context.Context, opts models.ListOptions) (int, error) {
	var count int

	if err := r.db.GetWithContext(ctx, &count, SelectClustersFullMView, opts, "count(*)"); err != nil {
		return count, fmt.Errorf("failed to list clusters: %w", err)
	}

	return count, nil
}

// GetClusterByID retrieves a cluster's details by its unique identifier.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice containing a single inventory.Cluster object.
// - An error if the query fails or the cluster ID does not exist.
func (r *clusterRepositoryImpl) GetClusterByID(ctx context.Context, clusterID string) (*db.ClusterDBResponse, error) {
	var cluster db.ClusterDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
		},
	}

	if err := r.db.GetWithContext(ctx, &cluster, SelectClustersFullMView, opts, "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &cluster, nil
}

// GetClusterAccountName retrieves the accountID associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A string representing the accountID.
// - An error if the query fails or the cluster ID does not exist.
func (r *clusterRepositoryImpl) GetClusterAccountName(ctx context.Context, clusterID string) (string, error) {
	cluster, err := r.GetClusterByID(ctx, clusterID)
	if err != nil {
		return "", err
	}

	return cluster.AccountID, nil
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
	cluster, err := r.GetClusterByID(ctx, clusterID)
	if err != nil {
		return "", err
	}

	return cluster.Region, nil
}

// GetClusterTags retrieves the tags associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice of inventory.Tag objects representing the cluster's tags.
// - An error if the query fails.
func (r *clusterRepositoryImpl) GetClusterTags(ctx context.Context, clusterID string) ([]db.TagDBResponse, error) {
	result := []db.TagDBResponse{}

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
		},
	}

	if err := r.db.SelectWithContext(ctx, &result, SelectClusterTags, opts, "cluster_id", "key, value"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return result, nil
}

// GetClustersOnAccount retrieves all clusters associated with a specific account.
//
// Parameters:
// - accountName: The ID of the account whose clusters will be retrieved.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (r *clusterRepositoryImpl) GetClustersOnAccount(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error) {
	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"account_id": accountID,
		},
	}

	clusters, _, err := r.ListClusters(ctx, opts)

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
func (r *clusterRepositoryImpl) GetInstancesOnCluster(ctx context.Context, clusterID string) ([]db.InstanceDBResponse, error) {
	instances := []db.InstanceDBResponse{}

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
		},
	}

	var cluster db.ClusterDBResponse
	if err := r.db.GetWithContext(ctx, &cluster, SelectClustersFullMView, opts, "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := r.db.SelectWithContext(ctx, &instances, SelectInstancesFullView, opts, "cluster_id", "*"); err != nil {
		return instances, fmt.Errorf("failed to list instances for cluster '%s': %w", clusterID, err)
	}

	return instances, nil
}

// GetClustersOverview returns a summary of cluster statuses
// It counts the number of clusters that are running, stopped or terminated.
func (r *clusterRepositoryImpl) GetClustersOverview(ctx context.Context) (inventory.ClustersSummary, error) {
	var countsDB inventory.ClustersSummary

	if err := r.db.GetWithContext(ctx, &countsDB, ClustersTable, models.ListOptions{},
		"COUNT(CASE WHEN status = 'Running' THEN 1 END) AS running",
		"COUNT(CASE WHEN status = 'Stopped' THEN 1 END) AS stopped",
		"COUNT(CASE WHEN status = 'Terminated' THEN 1 END) AS archived",
	); err != nil {
		return countsDB, fmt.Errorf("failed to list clusters: %w", err)
	}

	return countsDB, nil
}

// CreateClusters inserts a list of clusters into the database in a transaction.
//
// Parameters:
// - clusters: A slice of inventory.Cluster objects to insert.
//
// Returns:
// - An error if the transaction fails or the query encounters an issue.
func (r *clusterRepositoryImpl) CreateClusters(ctx context.Context, clusters []inventory.Cluster) error {
	if err := r.db.InsertWithContext(ctx, InsertClustersQuery, clusters); err != nil {
		return err
	}
	return nil
}

// UpdateCluster updates mutable fields of an existing cluster in the database.
// Only non-nil fields in the patch request will be updated.
func (r *clusterRepositoryImpl) UpdateCluster(ctx context.Context, clusterID string, patch dto.ClusterPatchRequest) (err error) {
	// Build dynamic UPDATE query with positional parameters
	query := "UPDATE clusters SET "
	args := make([]interface{}, 0)
	updateFields := make([]string, 0)
	argCount := 1

	if patch.ConsoleLink != nil {
		updateFields = append(updateFields, fmt.Sprintf("console_link = $%d", argCount))
		args = append(args, *patch.ConsoleLink)
		argCount++
	}

	if patch.Owner != nil {
		updateFields = append(updateFields, fmt.Sprintf("owner = $%d", argCount))
		args = append(args, *patch.Owner)
		argCount++
	}

	// If no fields to update, return early
	if len(updateFields) == 0 {
		return nil
	}

	query += strings.Join(updateFields, ", ")
	query += fmt.Sprintf(" WHERE cluster_id = $%d", argCount)
	args = append(args, clusterID)

	// Execute update in a transaction
	tx, txErr := r.db.NewTx(ctx)
	if txErr != nil {
		return txErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, execErr := tx.ExecContext(ctx, query, args...); execErr != nil {
		err = fmt.Errorf("exec UPDATE error: %w", execErr)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit UPDATE error: %w", err)
	}

	// Refresh materialized view after transaction commits
	// Note: REFRESH MATERIALIZED VIEW must run outside the transaction
	tx2, tx2Err := r.db.NewTx(ctx)
	if tx2Err != nil {
		return fmt.Errorf("failed to create transaction for refresh: %w", tx2Err)
	}
	defer func() {
		if err != nil {
			_ = tx2.Rollback()
		}
	}()

	if _, refreshErr := tx2.ExecContext(ctx, "REFRESH MATERIALIZED VIEW m_clusters_full_view"); refreshErr != nil {
		err = fmt.Errorf("refresh materialized view error: %w", refreshErr)
		return err
	}

	err = tx2.Commit()
	if err != nil {
		return fmt.Errorf("commit refresh error: %w", err)
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
func (r *clusterRepositoryImpl) UpdateClusterStatusByClusterID(_ context.Context, _ string, _ string) error {
	// TODO
	return nil
}

// DeleteCluster deletes a cluster from the database.
//
// Parameters:
// - id: The id of the cluster to delete.
//
// Returns:
// - An error if the database transaction fails.
func (r *clusterRepositoryImpl) DeleteCluster(ctx context.Context, clusterID string) error {
	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
		},
	}

	if err := r.db.DeleteWithContext(ctx, ClustersTable, opts); err != nil {
		return err
	}
	return nil
}
