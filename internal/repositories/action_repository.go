package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
)

const (
	// DB Table for actions
	ScheduleTable = "schedule"
	// View for getting scheduled actions
	SelectScheduleFullView = "schedule_full_view"
	// EnableActionQuery enables the action to be re-scheduled on next agent polling
	EnableActionQuery = `
		UPDATE
			schedule
		SET
			enabled = true
		WHERE id = $1
	`
	// DisableActionQuery disables the action to don't be re-scheduled on next agent polling
	DisableActionQuery = `
		UPDATE
			schedule
		SET
			enabled = false
		WHERE id = $1
	`
	// InsertAction inserts a new action returning the ID
	InsertActionsQuery = `
		INSERT INTO schedule (
			type,
			time,
			operation,
			target,
			status,
			enabled
		) VALUES (
			:type,
			(SELECT now()),
			:operation,
			(SELECT id FROM clusters WHERE cluster_id = :target.cluster_id),
			:status,
			:enabled
		) RETURNING id
	`
	// InsertScheduledActionQuery inserts new scheduled actions on the DB
	InsertScheduledActionsQuery = `
		INSERT INTO schedule (
			type,
			time,
			operation,
			target,
			status,
			enabled
		) VALUES (
			'scheduled_action',
			:time,
			:operation,
			:target.cluster_id,
			:status,
			:enabled
		)
	`
	// InsertCronActionQuery inserts new Cron actions on the DB
	InsertCronActionsQuery = `
		INSERT INTO schedule (
			type,
			cron_exp,
			operation,
			target,
			status,
			enabled
		) VALUES (
			'cron_action',
			:cron_exp,
			:operation,
			:target.cluster_id,
			:status,
			:enabled
		)
	`
	// UpdateActionQuery updates a single action on the DB
	UpdateActionQuery = `
		UPDATE schedule
		SET
			status = :status
		WHERE id = :id
	`
)

var _ ActionRepository = (*actionRepositoryImpl)(nil)

// ActionRepository defines the interface for data access operations for actions.
type ActionRepository interface {
	List(ctx context.Context, opts models.ListOptions) ([]db.ActionDBResponse, int, error)
	GetByID(ctx context.Context, actionID string) (db.ActionDBResponse, error)
	Enable(ctx context.Context, actionID string) error
	Disable(ctx context.Context, actionID string) error
	Create(ctx context.Context, newActions []actions.Action) error
	CreateAction(ctx context.Context, action actions.Action) (int64, error)
	Delete(ctx context.Context, actionID string) error
	Update(ctx context.Context, action actions.Action) error
}

type actionRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewActionRepository(db *dbclient.DBClient) ActionRepository {
	return &actionRepositoryImpl{db: db}
}

// Lists runs the db select query for retrieving the scheduled actions on the DB
//
// Parameters:
//
// Returns:
//   - An array of actions.Action with the scheduled actions declared on the DB
//   - An error if the query fails
func (r *actionRepositoryImpl) List(ctx context.Context, opts models.ListOptions) ([]db.ActionDBResponse, int, error) {
	var schedule []db.ActionDBResponse

	if err := r.db.SelectWithContext(ctx, &schedule, SelectScheduleFullView, opts, "id", "*"); err != nil {
		return schedule, 0, fmt.Errorf("failed to list schedule: %w", err)
	}

	return schedule, len(schedule), nil
}

// Enable enables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) Enable(ctx context.Context, actionID string) error {
	return r.db.UpdateWithContext(ctx, EnableActionQuery, actionID)
}

// Disable Disables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) Disable(ctx context.Context, actionID string) error {
	return r.db.UpdateWithContext(ctx, DisableActionQuery, actionID)
}

// GetByID runs the db select query for retrieving a specific scheduled action by its ID
//
// Parameters:
//
// Returns:
//   - An actions.Action object
//   - An error if the query fails
func (r *actionRepositoryImpl) GetByID(ctx context.Context, actionID string) (db.ActionDBResponse, error) {
	var action db.ActionDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"id": actionID,
		},
	}

	if err := r.db.GetWithContext(ctx, &action, SelectScheduleFullView, opts, "id", "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return action, ErrNotFound
		}
		return action, err
	}

	return action, nil
}

// Create creates a batch of scheduled actions in the database.
//
// Parameters:
//   - An array of actions.Action to write on the DB
//
// Returns:
//   - An error if the insert fails
//
// TODO: Temporal fix returning TX from DBClient to manage both insertions in the same sql transaction
func (r *actionRepositoryImpl) Create(ctx context.Context, newActions []actions.Action) error {
	schedActions, cronActions := actions.SplitActionsByType(newActions)

	tx, err := r.db.NewTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Writing Scheduled Actions
	if len(schedActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, InsertScheduledActionsQuery, schedActions); err != nil {
			return fmt.Errorf("failed to insert scheduled actions: %w", err)
		}
	}

	// Writing Cron Actions
	if len(cronActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, InsertCronActionsQuery, cronActions); err != nil {
			return fmt.Errorf("failed to insert cron actions: %w", err)
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// AddEvent inserts a new audit event into the database and returns the event ID.
func (r *actionRepositoryImpl) CreateAction(ctx context.Context, action actions.Action) (int64, error) {
	var returnedValue int64
	returnedValue, err := r.db.InsertWithReturnWithContext(ctx, InsertActionsQuery, action)
	if err != nil {
		return -1, err
	}

	return returnedValue, nil
}

// Delete removes an actions.ScheduledAction action from the DB based on its ID
//
// Parameters:
//   - A string containing the action ID to be removed
//
// Returns:
//   - An error if the delete query fails
func (r *actionRepositoryImpl) Delete(ctx context.Context, actionID string) error {
	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"id": actionID,
		},
	}

	if err := r.db.DeleteWithContext(ctx, ScheduleTable, opts); err != nil {
		return err
	}
	return nil
}

// Update updates an actions.Action action from the DB based on its ID
//
// Parameters:
//   - An action with an exisiting ID to be updated
//
// Returns:
//   - An error if the delete query fails
func (r *actionRepositoryImpl) Update(ctx context.Context, action actions.Action) error {
	if _, err := r.GetByID(ctx, action.GetID()); err != nil {
		return err
	}

	if err := r.db.NamedUpdateWithContext(ctx, UpdateActionQuery, action); err != nil {
		return err
	}

	return nil
}
