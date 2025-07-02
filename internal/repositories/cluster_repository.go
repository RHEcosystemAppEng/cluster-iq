package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/jmoiron/sqlx"
)

// HUGE TODO!!! Update the entire code to use CONTEXT

var ErrNotFound = errors.New("requested resource not found")

var _ ClusterRepository = (*ClusterRepositoryImpl)(nil)

type ClusterRepository interface {
	GetClusterByID(clusterID string) (inventory.Cluster, error)
	ListClusters(ctx context.Context, opts ListOptions) ([]inventory.Cluster, int, error)
	// TODO: review if need to move
	GetClustersOverview() (models.ClustersSummary, error)
	GetClusterAccountName(clusterID string) (string, error)
	GetClusterRegion(clusterID string) (string, error)
	GetClusterTags(clusterID string) ([]inventory.Tag, error)
	GetClustersOnAccount(accountName string) ([]inventory.Cluster, error)
	GetInstancesOnCluster(clusterID string) ([]inventory.Instance, error)
	DeleteCluster(id string) error
}
type ClusterRepositoryImpl struct {
	db *sqlx.DB
}

func NewClusterRepository(db *sqlx.DB) *ClusterRepositoryImpl {
	return &ClusterRepositoryImpl{db: db}
}

// GetClusterByID retrieves a cluster's details by its unique identifier.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice containing a single inventory.Cluster object.
// - An error if the query fails or the cluster ID does not exist.
func (r *ClusterRepositoryImpl) GetClusterByID(clusterID string) (inventory.Cluster, error) {
	var cluster inventory.Cluster
	err := r.db.Get(&cluster, SelectClustersByIDQuery, clusterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return inventory.Cluster{}, ErrNotFound
		}
		return inventory.Cluster{}, err
	}
	return cluster, nil
}

// ListClusters retrieves all clusters from the database.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (r *ClusterRepositoryImpl) ListClusters(ctx context.Context, opts ListOptions) ([]inventory.Cluster, int, error) {
	baseQuery := SelectClustersQuery
	countQuery := "SELECT COUNT(*) FROM clusters"

	whereClauses, namedArgs := buildWhereClauses(opts.Filters)
	if len(whereClauses) > 0 {
		whereStr := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += whereStr
		countQuery += whereStr
	}

	// Counting
	countStmt, countArgs, err := sqlx.Named(countQuery, namedArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to bind named args for count query: %w", err)
	}
	countStmt = r.db.Rebind(countStmt)

	var total int
	if err := r.db.GetContext(ctx, &total, countStmt, countArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to execute count query: %w", err)
	}

	// Pagination
	baseQuery += " LIMIT :pagesize OFFSET :offset"
	namedArgs["pagesize"] = opts.PageSize
	namedArgs["offset"] = opts.Offset

	// Main select
	queryStmt, queryArgs, err := sqlx.Named(baseQuery, namedArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to bind named args for select query: %w", err)
	}
	queryStmt = r.db.Rebind(queryStmt)

	var clusters []inventory.Cluster
	if err := r.db.SelectContext(ctx, &clusters, queryStmt, queryArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to execute select query: %w", err)
	}

	return clusters, total, nil
}

// GetClustersOverview returns a summary of cluster statuses
// It counts the number of clusters that are running, stopped or terminated.
func (r *ClusterRepositoryImpl) GetClustersOverview() (models.ClustersSummary, error) {
	var clustersOverview models.ClustersSummary
	err := r.db.Get(&clustersOverview, SelectClustersOverview)
	return clustersOverview, err
}

// GetClusterAccountName retrieves the account name associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A string representing the account name.
// - An error if the query fails or the cluster ID does not exist.
func (r *ClusterRepositoryImpl) GetClusterAccountName(clusterID string) (string, error) {
	var accountName string
	err := r.db.Get(&accountName, SelectClusterAccountNameQuery, clusterID)
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
func (r *ClusterRepositoryImpl) GetClusterRegion(clusterID string) (string, error) {
	var region string
	err := r.db.Get(&region, SelectClusterRegionQuery, clusterID)
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
func (r *ClusterRepositoryImpl) GetClusterTags(clusterID string) ([]inventory.Tag, error) {
	var tags []inventory.Tag
	err := r.db.Select(&tags, SelectClusterTags, clusterID)
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
func (r *ClusterRepositoryImpl) GetClustersOnAccount(accountName string) ([]inventory.Cluster, error) {
	var clusters []inventory.Cluster
	err := r.db.Select(&clusters, SelectClustersOnAccountQuery, accountName)
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
func (r *ClusterRepositoryImpl) GetInstancesOnCluster(clusterID string) ([]inventory.Instance, error) {
	var instances []inventory.Instance
	err := r.db.Select(&instances, SelectInstancesOnClusterQuery, clusterID)
	return instances, err
}

// DeleteCluster deletes a cluster from the database.
//
// Parameters:
// - clusterName: The name of the cluster to delete.
//
// Returns:
// - An error if the database transaction fails.
func (r *ClusterRepositoryImpl) DeleteCluster(clusterName string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// TODO. Review
				// r.logger.Error("Failed to rollback DeleteCluster transaction", zap.Error(rbErr))
			}
		}
	}()

	tx.MustExec(DeleteClusterQuery, clusterName)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func buildWhereClauses(filters map[string]interface{}) ([]string, map[string]interface{}) {
	args := make(map[string]interface{})
	if len(filters) == 0 {
		return nil, args
	}

	allowedFilters := map[string]bool{
		"status":       true,
		"account_name": true,
		"provider":     true,
		"region":       true,
	}

	clauses := make([]string, 0, len(filters))

	for key, value := range filters {
		if allowedFilters[key] {
			clauses = append(clauses, fmt.Sprintf("%s = :%s", key, key))
			args[key] = value
		}
	}
	return clauses, args
}
