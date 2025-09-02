package sqlclient

import (
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	dbmodel "github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"go.uber.org/zap"
)

const (
	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = `
		SELECT
			*
		FROM instances_full_view_with_tags
		ORDER BY instance_id
	`

	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = `
		SELECT
			*
		FROM instances_full_view_with_tags
		WHERE instance_id = $1
		ORDER BY instance_id
	`

	// SelectInstancesOverview returns the total count of all instances
	SelectInstancesOverview = `
		SELECT
			COUNT(CASE WHEN status = 'Running' THEN 1 END) AS running,
			COUNT(CASE WHEN status = 'Stopped' THEN 1 END) AS stopped,
			COUNT(CASE WHEN status = 'Terminated' THEN 1 END) AS archived
		FROM instances;
	`

	// InsertInstancesQuery inserts into a new instance in its table
	InsertInstancesQuery = `
		INSERT INTO instances (
			instance_id,
			instance_name,
			instance_type,
			provider,
			availability_zone,
			status,
			cluster_id,
			last_scan_ts,
			created_at,
			age
		) VALUES (
			:instance_id,
			:instance_name,
			:instance_type,
			:provider,
			:availability_zone,
			:status,
			:cluster_id,
			:last_scan_ts,
			:created_at,
			:age
		) ON CONFLICT (id) DO UPDATE SET
			instance_name = EXCLUDED.instance_name,
			instance_type = EXCLUDED.instance_type,
			provider = EXCLUDED.provider,
			availability_zone = EXCLUDED.availability_zone,
			status = EXCLUDED.status,
			cluster_id = EXCLUDED.cluster_id,
			last_scan_ts = EXCLUDED.last_scan_ts,
			created_at = EXCLUDED.created_at,
			age = EXCLUDED.age
	`

	// DeleteInstanceQuery removes an instance by its ID
	DeleteInstanceQuery = `DELETE FROM instances WHERE instance_id = $1`
)

// GetInstances retrieves all instances from the database and maps them to inventory.Instance objects.
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (a SQLClient) GetInstances() ([]dbmodel.InstanceDBResponse, error) {
	var dbinstances []dbmodel.InstanceDBResponse
	if err := a.db.Select(&dbinstances, SelectInstancesQuery); err != nil {
		return nil, err
	}

	return dbinstances, nil
}

// GetInstanceByID retrieves an instance by its ID.
//
// Parameters:
// - instanceID: The ID of the instance to retrieve.
//
// Returns:
// - A slice of inventory.Instance objects (usually one element).
// - An error if the query fails.
func (a SQLClient) GetInstanceByID(instanceID string) (*dbmodel.InstanceDBResponse, error) {
	var dbinstances dbmodel.InstanceDBResponse
	if err := a.db.Get(&dbinstances, SelectInstancesByIDQuery, instanceID); err != nil {
		return nil, err
	}

	return &dbinstances, nil
}

func (a SQLClient) GetClusterInternalID(clusterID string) (int, error) {
	var id int
	if err := a.db.Get(&id, fmt.Sprintf("SELECT id FROM clusters WHERE cluster_id = '%s'", clusterID)); err != nil {
		return -1, err
	}
	return id, nil
}

// WriteInstances writes a batch of instances and their tags to the database in a transaction.
//
// Parameters:
// - instances: A slice of inventory.Instance objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (a SQLClient) WriteInstances(instances []inventory.Instance) error {
	var tags []inventory.Tag
	for _, instance := range instances {
		tags = append(tags, instance.Tags...)
	}

	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback WriteInstances query", zap.Error(rbErr))
			}
		}
	}()
	// Writing Instances
	if _, err := tx.NamedExec(InsertInstancesQuery, instances); err != nil {
		a.logger.Error("Failed to prepare InsertInstancesQuery query", zap.Error(err))
		return err
	}

	// Writing tags
	if len(tags) > 0 {
		if _, err := tx.NamedExec(InsertTagsQuery, tags); err != nil {
			a.logger.Error("Failed to prepare InsertTagsQuery query", zap.Error(err))
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// DeleteInstance deletes an instance and its associated tags from the database.
//
// Parameters:
// - instanceID: The ID of the instance to delete.
//
// Returns:
// - An error if the transaction fails.
func (a SQLClient) DeleteInstance(instanceID string) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback DeleteInstance transaction", zap.Error(rbErr))
			}
		}
	}()

	// TODO review if ON CASCADE delete works
	//tx.MustExec(DeleteTagsQuery, instanceID)
	tx.MustExec(DeleteInstanceQuery, instanceID)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
