package sqlclient

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"go.uber.org/zap"
)

const (
	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = `
		SELECT * FROM instances
		JOIN tags ON
			instances.id = tags.instance_id
		WHERE id = $1
		ORDER BY name
	`

	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = `
		SELECT * FROM instances
		JOIN tags ON
			instances.id = tags.instance_id
		ORDER BY name
	`

	// InsertInstancesQuery inserts into a new instance in its table
	InsertInstancesQuery = `
		INSERT INTO instances (
			id,
			name,
			provider,
			instance_type,
			availability_zone,
			status,
			cluster_id,
			last_scan_ts,
			created_at,
			age,
			daily_cost,
			total_cost
		) VALUES (
			:id,
			:name,
			:provider,
			:instance_type,
			:availability_zone,
			:status,
			:cluster_id,
			:last_scan_ts,
			:created_at,
			:age,
			:daily_cost,
			:total_cost
		) ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			provider = EXCLUDED.provider,
			instance_type = EXCLUDED.instance_type,
			availability_zone = EXCLUDED.availability_zone,
			status = EXCLUDED.status,
			cluster_id = EXCLUDED.cluster_id,
			last_scan_ts = EXCLUDED.last_scan_ts,
			created_at = EXCLUDED.created_at,
			age = EXCLUDED.age
	`

	// DeleteInstanceQuery removes an instance by its ID
	DeleteInstanceQuery = `DELETE FROM instances WHERE id=$1`
)

// GetInstances retrieves all instances from the database and maps them to inventory.Instance objects.
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (a SQLClient) GetInstances() ([]models.InstanceDBResponse, error) {
	var dbinstances []models.InstanceDBResponse
	if err := a.db.Select(&dbinstances, SelectInstancesQuery); err != nil {
		return nil, err
	}

	instances := joinInstancesTags(dbinstances)

	return instances, nil
}

// GetInstanceByID retrieves an instance by its ID.
//
// Parameters:
// - instanceID: The ID of the instance to retrieve.
//
// Returns:
// - A slice of inventory.Instance objects (usually one element).
// - An error if the query fails.
func (a SQLClient) GetInstanceByID(instanceID string) ([]inventory.Instance, error) {
	var dbinstances []models.InstanceDB
	if err := a.db.Select(&dbinstances, SelectInstancesByIDQuery, instanceID); err != nil {
		return nil, err
	}

	instances := joinInstancesTags(dbinstances)

	return instances, nil
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
	if _, err := tx.NamedExec(InsertTagsQuery, tags); err != nil {
		a.logger.Error("Failed to prepare InsertTagsQuery query", zap.Error(err))
		return err
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

	tx.MustExec(DeleteTagsQuery, instanceID)
	tx.MustExec(DeleteInstanceQuery, instanceID)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// joinInstancesTags maps an array of InstanceDBResponse objects into a slice of inventory.Instance objects.
//
// Parameters:
// - dbinstances: A slice of InstanceDB objects.
//
// Returns:
// - A slice of inventory.Instance objects.
func joinInstancesTags(dbinstances []models.InstanceDBResponse) []inventory.Instance {
	instanceMap := make(map[string]*inventory.Instance)
	for _, dbinstance := range dbinstances {
		if _, ok := instanceMap[dbinstance.ID]; ok {
			// Adding tag to an already read instance
			instance := instanceMap[dbinstance.ID]
			instance.AddTag(
				*inventory.NewTag(dbinstance.TagKey, dbinstance.TagValue, dbinstance.ID),
			)
		} else {
			// Adding a new instance to the response
			instanceMap[dbinstance.ID] = inventory.NewInstance(
				dbinstance.ID,
				dbinstance.Name,
				dbinstance.Provider,
				dbinstance.InstanceType,
				dbinstance.AvailabilityZone,
				dbinstance.Status,
				[]inventory.Tag{*inventory.NewTag(dbinstance.TagKey, dbinstance.TagValue, dbinstance.ID)},
				dbinstance.CreationTimestamp,
			)
			instanceMap[dbinstance.ID].ClusterID = dbinstance.ClusterID
		}
	}

	// Converting map into list
	var instances []inventory.Instance
	for _, instance := range instanceMap {
		instances = append(instances, *instance)
	}

	return instances
}
