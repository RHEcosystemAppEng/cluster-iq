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

var _ InstanceRepository = (*instanceRepositoryImpl)(nil)

// InstanceRepository defines the interface for data access operations for instances.
type InstanceRepository interface {
	ListInstances(ctx context.Context, opts models.ListOptions) ([]db.InstanceDBResponse, int, error)
	GetInstanceByID(ctx context.Context, instanceID string) (db.InstanceDBResponse, error)
	GetInstancesOutdatedBilling(ctx context.Context) ([]db.InstanceDBResponse, error)
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

	if err := r.db.Select(&instances, SelectInstancesFullView, opts, "instance_id", "*"); err != nil {
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

	if err := r.db.Get(&instance, SelectInstancesFullWithTagsView, opts, "instance_id", "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return instance, ErrNotFound
		}
		return instance, err
	}
	return instance, nil
}

// GetInstancesOutdatedBilling retrieves instances with outdated billing information.
//
// Parameters:
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (r *instanceRepositoryImpl) GetInstancesOutdatedBilling(ctx context.Context) ([]db.InstanceDBResponse, error) {
	var instances []db.InstanceDBResponse

	if err := r.db.Select(instances, SelectInstancesPendingExpenseUpdate, models.ListOptions{}, "instance_id", "*"); err != nil {
		return instances, fmt.Errorf("failed to list instances pending of expense update: %w", err)
	}

	return instances, nil
}

// GetInstancesOverview returns a summary of instances grouped by their status.
// It provides the total count along with counts of running and stopped instances.
func (r *instanceRepositoryImpl) GetInstancesOverview(ctx context.Context) (inventory.InstancesSummary, error) {
	var countsDB inventory.InstancesSummary

	if err := r.db.Select(&countsDB, InstancesTable, models.ListOptions{}, "",
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
	fmt.Println("INSERTING INSTANCES", instances)

	if err := r.db.Insert(InsertInstancesQuery, instances); err != nil {
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
		if err := r.db.Insert(InsertTagsQuery, newTags); err != nil {
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

	if err := r.db.Delete(InstancesTable, opts); err != nil {
		return err
	}
	return nil
}

func (r *instanceRepositoryImpl) parseInstanceInternalID(ctx context.Context, instance inventory.Instance) (string, error) {
	var id string

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"instance_id": instance.InstanceID,
		},
	}

	if err := r.db.Get(&id, AccountsTable, opts, "id"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	return id, nil
}
