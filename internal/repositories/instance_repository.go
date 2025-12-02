package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
)

const (
	// DB Table for instances
	InstancesTable = "instances"
	// View for SELECT operations on Instances
	SelectInstancesFullView = "instances_full_view"
	// View for SELECT operations on Instances
	SelectInstancesFullMView = "m_instances_full_view"
	// View for SELECT operations on Instances including tags
	SelectInstancesFullWithTagsView = "instances_full_view_with_tags"
	// View for SELECT operations on Instances including tags
	SelectInstancesFullWithTagsMView = "m_instances_full_view_with_tags"
	// InsertInstancesQuery inserts into a new instance in its table
	InsertInstancesQuery = `
		INSERT INTO instances (
			instance_id,
			instance_name,
			provider,
			instance_type,
			availability_zone,
			status,
			cluster_id,
			last_scan_ts,
			created_at,
			age
		) VALUES (
			:instance_id,
			:instance_name,
			:provider,
			:instance_type,
			:availability_zone,
			:status,
			(SELECT id FROM clusters WHERE cluster_id = :cluster_id),
			:last_scan_ts,
			:created_at,
			:age
		) ON CONFLICT (instance_id, cluster_id) DO UPDATE SET
			instance_name = EXCLUDED.instance_name,
			provider = EXCLUDED.provider,
			instance_type = EXCLUDED.instance_type,
			availability_zone = EXCLUDED.availability_zone,
			status = EXCLUDED.status,
			last_scan_ts = EXCLUDED.last_scan_ts,
			created_at = EXCLUDED.created_at,
			age = EXCLUDED.age
	`
	// InsertTagsQuery inserts into a new tag for an instance
	InsertTagsQuery = `
		INSERT INTO tags (
			key,
			value,
			instance_id
		) VALUES (
			:key,
			:value,
			(SELECT id FROM instances WHERE instance_id = :instance_id)
		) ON CONFLICT (key, instance_id) DO UPDATE SET
			value = EXCLUDED.value
	`
)

var _ InstanceRepository = (*instanceRepositoryImpl)(nil)

// InstanceRepository defines the interface for data access operations for instances.
type InstanceRepository interface {
	ListInstances(ctx context.Context, opts models.ListOptions) ([]db.InstanceDBResponse, int, error)
	GetInstanceByID(ctx context.Context, instanceID string) (db.InstanceDBResponse, error)
	GetInstancesOverview(ctx context.Context) (inventory.InstancesSummary, error)
	CreateInstances(ctx context.Context, instances []inventory.Instance) error
	DeleteInstance(ctx context.Context, instanceID string) error
}

type instanceRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewInstanceRepository(db *dbclient.DBClient) InstanceRepository {
	return &instanceRepositoryImpl{db: db}
}

// ListInstances retrieves all instances from the database and maps them to inventory.Instance objects.
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (r *instanceRepositoryImpl) ListInstances(ctx context.Context, opts models.ListOptions) ([]db.InstanceDBResponse, int, error) {
	var instances []db.InstanceDBResponse

	if err := r.db.SelectWithContext(ctx, &instances, SelectInstancesFullMView, opts, "instance_id", "instance_id"); err != nil {
		return instances, 0, fmt.Errorf("failed to list instances: %w", err)
	}

	return instances, len(instances), nil
}

// GetInstanceByID retrieves an instance by its ID.
//
// Parameters:
// - instanceID: The ID of the instance to retrieve.
//
// Returns:
// - An inventory.Instance object.
// - An error if the query fails.
func (r *instanceRepositoryImpl) GetInstanceByID(ctx context.Context, instanceID string) (db.InstanceDBResponse, error) {
	var instance db.InstanceDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"instance_id": instanceID,
		},
	}

	if err := r.db.GetWithContext(ctx, &instance, SelectInstancesFullWithTagsMView, opts, "instance_id", "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return instance, ErrNotFound
		}
		return instance, err
	}
	return instance, nil
}

// GetInstancesOverview returns a summary of instances grouped by their status.
// It provides the total count along with counts of running and stopped instances.
func (r *instanceRepositoryImpl) GetInstancesOverview(ctx context.Context) (inventory.InstancesSummary, error) {
	var countsDB inventory.InstancesSummary

	if err := r.db.GetWithContext(ctx, &countsDB, InstancesTable, models.ListOptions{},
		"COUNT(CASE WHEN status = 'Running' THEN 1 END) AS running",
		"COUNT(CASE WHEN status = 'Stopped' THEN 1 END) AS stopped",
		"COUNT(CASE WHEN status = 'Terminated' THEN 1 END) AS archived",
	); err != nil {
		return countsDB, fmt.Errorf("failed to list clusters: %w", err)
	}

	return countsDB, nil
}

// CreateInstances writes a batch of instances and their tags to the database in a transaction.
//
// Parameters:
// - instances: A slice of inventory.Instance objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (r *instanceRepositoryImpl) CreateInstances(ctx context.Context, instances []inventory.Instance) error {
	if err := r.db.InsertWithContext(ctx, InsertInstancesQuery, instances); err != nil {
		return err
	}

	var newTags []inventory.Tag
	for _, instance := range instances {
		for _, tag := range instance.Tags {
			tag.InstanceID = instance.InstanceID
			newTags = append(newTags, tag)
		}
	}

	if len(newTags) > 0 {
		if err := r.db.InsertWithContext(ctx, InsertTagsQuery, newTags); err != nil {
			return err
		}
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
func (r *instanceRepositoryImpl) DeleteInstance(ctx context.Context, instanceID string) error {
	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"instance_id": instanceID,
		},
	}

	if err := r.db.DeleteWithContext(ctx, InstancesTable, opts); err != nil {
		return err
	}
	return nil
}
