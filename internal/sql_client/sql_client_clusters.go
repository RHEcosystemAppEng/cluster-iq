package sqlclient

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"go.uber.org/zap"
)

const (
	// SelectClustersQuery returns every cluster in the inventory ordered by Name
	SelectClustersQuery = `
		SELECT * FROM clusters_full_view
		ORDER BY cluster_name
	`

	// SelectClustersByIDuery returns an cluster by its Name
	SelectClustersByIDuery = `
		SELECT * FROM clusters_full_view
		WHERE cluster_id = $1
		ORDER BY cluster_name
	`

	// SelectInstancesOnClusterQuery returns every instance belonging to a cluster
	SelectInstancesOnClusterQuery = `
		SELECT * FROM instances_full_view
		WHERE cluster_id = $1
		ORDER BY instance_name
	`

	// InsertClustersQuery inserts into a new instance in its table
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
			:account_id,
			:console_link,
			:last_scan_ts,
			:created_at,
			:age,
			:owner
		) ON CONFLICT (id) DO UPDATE SET
			provider = EXCLUDED.provider,
			status = EXCLUDED.status,
			region = EXCLUDED.region,
			console_link = EXCLUDED.console_link,
			last_scan_ts = EXCLUDED.last_scan_ts,
			created_at = EXCLUDED.created_at,
			age = EXCLUDED.age,
			owner = EXCLUDED.owner
	`

	// DeleteClusterQuery removes an cluster by its name
	DeleteClusterQuery = `DELETE FROM clusters WHERE cluster_id=$1`
)

// GetClusters retrieves all clusters from the database.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (a SQLClient) GetClusters() ([]models.ClusterDBResponse, error) {
	var clusters []models.ClusterDBResponse
	if err := a.db.Select(&clusters, SelectClustersQuery); err != nil {
		return nil, err
	}
	return clusters, nil
}

// GetClusterByID retrieves a cluster's details by its unique identifier.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice containing a single inventory.Cluster object.
// - An error if the query fails or the cluster ID does not exist.
func (a SQLClient) GetClusterByID(clusterID string) (*models.ClusterDBResponse, error) {
	var cluster models.ClusterDBResponse
	if err := a.db.Get(&cluster, SelectClustersByIDuery, clusterID); err != nil {
		return nil, err
	}
	return &cluster, nil
}

// GetClusterTags retrieves the tags associated with a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice of inventory.Tag objects representing the cluster's tags.
// - An error if the query fails.
func (a SQLClient) GetClusterTags(clusterID string) ([]inventory.Tag, error) {
	var tags []inventory.Tag
	if err := a.db.Select(&tags, SelectClusterTags, clusterID); err != nil {
		return nil, err
	}
	return tags, nil
}

// GetInstancesOnCluster retrieves all instances belonging to a specific cluster.
//
// Parameters:
// - clusterID: The unique identifier of the cluster.
//
// Returns:
// - A slice of inventory.Instance objects representing the instances in the cluster.
// - An error if the query fails.
func (a SQLClient) GetInstancesOnCluster(clusterID string) ([]models.InstanceDBResponse, error) {
	var instances []models.InstanceDBResponse
	if err := a.db.Select(&instances, SelectInstancesOnClusterQuery, clusterID); err != nil {
		return nil, err
	}
	return instances, nil
}

func (a SQLClient) GetAccountInternalID(accountID string) (int, error) {
	var id int
	if err := a.db.Get(&id, fmt.Sprintf("SELECT id FROM accounts WHERE account_id = '%s'", accountID)); err != nil {
		return -1, err
	}
	return id, nil
}

// WriteClusters inserts a list of clusters into the database in a transaction.
//
// Parameters:
// - clusters: A slice of inventory.Cluster objects to insert.
//
// Returns:
// - An error if the transaction fails or the query encounters an issue.
func (a SQLClient) WriteClusters(clusters []inventory.Cluster) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback WriteClusters transaction", zap.Error(rbErr))
			}
		}
	}()

	result, err := tx.NamedExec(InsertClustersQuery, clusters)
	if err != nil {
		if result != nil {
			rows, _ := result.RowsAffected()
			a.logger.Error("Failed to run InsertClustersQuery query", zap.Int64("rows_affected", rows))
		}
		a.logger.Error("Failed to prepare InsertClustersQuery query", zap.Error(err))
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// DeleteCluster deletes a cluster from the database.
//
// Parameters:
// - clusterID: The name of the cluster to delete.
//
// Returns:
// - An error if the database transaction fails.
func (a SQLClient) DeleteCluster(clusterID string) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback DeleteCluster transaction", zap.Error(rbErr))
			}
		}
	}()

	tx.MustExec(DeleteClusterQuery, clusterID)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
