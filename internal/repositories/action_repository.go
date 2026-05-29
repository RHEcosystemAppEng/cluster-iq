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

	InsertTargetQuery = `INSERT INTO targets (target_type, select_all) VALUES ($1, $2) RETURNING id`

	LinkTargetClusterQuery = `INSERT INTO target_clusters (target_id, cluster_id) SELECT $1, id FROM clusters WHERE cluster_id = $2`

	LinkTargetAccountQuery = `INSERT INTO target_accounts (target_id, account_id) SELECT $1, id FROM accounts WHERE account_id = $2`

	InsertScheduledActionWithTargetQuery = `
		INSERT INTO schedule (type, time, operation, target, status, enabled)
		VALUES ('scheduled_action', $1, $2, $3, $4, $5)
		RETURNING id
	`

	InsertCronActionWithTargetQuery = `
		INSERT INTO schedule (type, cron_exp, operation, target, status, enabled)
		VALUES ('cron_action', $1, $2, $3, $4, $5)
		RETURNING id
	`

	InsertInstantActionWithTargetQuery = `
		INSERT INTO schedule (type, time, operation, target, status, enabled)
		VALUES ('instant_action', NOW(), $1, $2, $3, $4)
		RETURNING id
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
	schedule := []db.ActionDBResponse{}

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
func (r *actionRepositoryImpl) Create(ctx context.Context, newActions []actions.Action) (err error) {
	tx, err := r.db.NewTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, action := range newActions {
		targetID, targetErr := createTargetForAction(ctx, tx, action)
		if targetErr != nil {
			return fmt.Errorf("failed to create target: %w", targetErr)
		}

		switch a := action.(type) {
		case *actions.ScheduledAction:
			_, err = tx.ExecContext(ctx, InsertScheduledActionWithTargetQuery,
				a.When, a.Operation, targetID, a.Status, a.Enabled)
		case *actions.CronAction:
			_, err = tx.ExecContext(ctx, InsertCronActionWithTargetQuery,
				a.Expression, a.Operation, targetID, a.Status, a.Enabled)
		case *actions.InstantAction:
			_, err = tx.ExecContext(ctx, InsertInstantActionWithTargetQuery,
				a.Operation, targetID, a.Status, a.Enabled)
		default:
			return fmt.Errorf("unsupported action type for batch create: %T", action)
		}
		if err != nil {
			return fmt.Errorf("failed to insert schedule: %w", err)
		}
	}

	return tx.Commit()
}

func (r *actionRepositoryImpl) CreateAction(ctx context.Context, action actions.Action) (int64, error) {
	tx, err := r.db.NewTx(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	targetID, err := createTargetForAction(ctx, tx, action)
	if err != nil {
		return -1, fmt.Errorf("failed to create target: %w", err)
	}

	var scheduleID int64
	err = tx.QueryRowContext(ctx, InsertInstantActionWithTargetQuery,
		action.GetActionOperation(), targetID, action.(*actions.InstantAction).Status, action.(*actions.InstantAction).Enabled,
	).Scan(&scheduleID)
	if err != nil {
		return -1, fmt.Errorf("failed to insert action: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return -1, fmt.Errorf("failed to commit action: %w", err)
	}

	return scheduleID, nil
}

func createTargetForAction(ctx context.Context, tx interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}, action actions.Action) (int64, error) {
	target := action.GetTarget()

	targetType := target.TargetType
	if targetType == "" {
		targetType = "Cluster"
	}

	var targetID int64
	if err := tx.QueryRowContext(ctx, InsertTargetQuery, targetType, target.SelectAll).Scan(&targetID); err != nil {
		return 0, fmt.Errorf("failed to insert target: %w", err)
	}

	switch targetType {
	case "Cluster":
		if target.ClusterID != "" {
			if _, err := tx.ExecContext(ctx, LinkTargetClusterQuery, targetID, target.ClusterID); err != nil {
				return 0, fmt.Errorf("failed to link target cluster: %w", err)
			}
		}
	case "Account":
		for _, accountID := range target.TargetAccountIDs {
			if _, err := tx.ExecContext(ctx, LinkTargetAccountQuery, targetID, accountID); err != nil {
				return 0, fmt.Errorf("failed to link target account: %w", err)
			}
		}
	}

	return targetID, nil
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
		// TODO include err message explaining the action for update doesn't exist
		return err
	}

	if err := r.db.NamedUpdateWithContext(ctx, UpdateActionQuery, action); err != nil {
		return err
	}

	return nil
}
