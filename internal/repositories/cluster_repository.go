package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	dbmodels "github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

const (
	// DB Table for clusters
	ClustersTable = "clusters"
	// View for getting clusters
	SelectClustersFullView = "clusters_full_view"
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
		) ON CONFLICT (id) DO UPDATE SET
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
	GetClusterByID(ctx context.Context, clusterID string) (*db.ClusterDBResponse, error)
	GetClusterAccountName(ctx context.Context, clusterID string) (string, error)
	GetClusterRegion(ctx context.Context, clusterID string) (string, error)
	GetClusterTags(ctx context.Context, clusterID string) ([]db.TagDBResponse, error)
	GetClustersOnAccount(ctx context.Context, accountName string) ([]db.ClusterDBResponse, error)
	GetInstancesOnCluster(ctx context.Context, clusterID string) ([]db.InstanceDBResponse, error)
	GetClustersOverview(ctx context.Context) (inventory.ClustersSummary, error)
	CreateClusters(ctx context.Context, clusters []inventory.Cluster) error
	UpdateCluster(ctx context.Context, cluster dto.ClusterDTORequest) error
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
	var clusters []dbmodels.ClusterDBResponse

	if err := r.db.Select(&clusters, SelectClustersFullMView, opts, "cluster_id", "*"); err != nil {
		return clusters, 0, fmt.Errorf("failed to list clusters: %w", err)
	}

	return clusters, len(clusters), nil
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
	var cluster dbmodels.ClusterDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
		},
	}

	if err := r.db.Get(&cluster, SelectClustersFullMView, opts, "*"); err != nil {
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
	var rawTags json.RawMessage
	var result []db.TagDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
		},
	}

	if err := r.db.Select(rawTags, SelectInstancesFullWithTagsMView, opts, "cluster_id", "DISTINCT ON (tags_json) tags_json"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, ErrNotFound
		}
		return result, err
	}

	var tags []inventory.Tag
	if err := json.Unmarshal(rawTags, &tags); err != nil {
		return result, err
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
	var instances []dbmodels.InstanceDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"cluster_id": clusterID,
		},
	}

	var cluster dbmodels.ClusterDBResponse
	if err := r.db.Get(&cluster, SelectClustersFullMView, opts, "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := r.db.Select(&instances, SelectInstancesFullView, opts, "cluster_id", "*"); err != nil {
		return instances, fmt.Errorf("failed to list instances for cluster '%s': %w", clusterID, err)
	}

	return instances, nil
}

// GetClustersOverview returns a summary of cluster statuses
// It counts the number of clusters that are running, stopped or terminated.
func (r *clusterRepositoryImpl) GetClustersOverview(ctx context.Context) (inventory.ClustersSummary, error) {
	var countsDB inventory.ClustersSummary

	if err := r.db.Get(&countsDB, ClustersTable, models.ListOptions{},
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
	if err := r.db.Insert(InsertClustersQuery, clusters); err != nil {
		return err
	}
	return nil
}

// UpdateCluster updates an existing cluster's details in the database.
func (r *clusterRepositoryImpl) UpdateCluster(ctx context.Context, cluster dto.ClusterDTORequest) error {
	//TODO
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
	//TODO
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

	if err := r.db.Delete(ClustersTable, opts); err != nil {
		return err
	}
	return nil
}
