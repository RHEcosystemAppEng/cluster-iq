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

var _ ActionRepository = (*actionRepositoryImpl)(nil)

// ActionRepository defines the interface for data access operations for actions.
type ActionRepository interface {
	ListScheduledActions(ctx context.Context, opts models.ListOptions) ([]db.ActionDBResponse, int, error)
	GetScheduledActionByID(ctx context.Context, actionID string) (db.ActionDBResponse, error)
	EnableScheduledAction(ctx context.Context, actionID string) error
	DisableScheduledAction(ctx context.Context, actionID string) error
	Create(ctx context.Context, newActions []actions.Action) error
	DeleteScheduledAction(ctx context.Context, actionID string) error
}

type actionRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewActionRepository(db *dbclient.DBClient) ActionRepository {
	return &actionRepositoryImpl{db: db}
}

// ListScheduledActions runs the db select query for retrieving the scheduled actions on the DB
//
// Parameters:
//
// Returns:
//   - An array of actions.Action with the scheduled actions declared on the DB
//   - An error if the query fails
func (r *actionRepositoryImpl) ListScheduledActions(ctx context.Context, opts models.ListOptions) ([]db.ActionDBResponse, int, error) {
	var schedule []db.ActionDBResponse

	if err := r.db.Select(schedule, "schedule_full_view", opts, "id", "*"); err != nil {
		return schedule, 0, fmt.Errorf("failed to list schedule: %w", err)
	}

	return schedule, len(schedule), nil
}

// EnableScheduledAction enables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) EnableScheduledAction(ctx context.Context, actionID string) error {
	return r.db.Update(EnableActionQuery, actionID)
}

// DisableScheduledAction Disables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) DisableScheduledAction(ctx context.Context, actionID string) error {
	return r.db.Update(DisableActionQuery, actionID)
}

// GetScheduledActionByID runs the db select query for retrieving a specific scheduled action by its ID
//
// Parameters:
//
// Returns:
//   - An actions.Action object
//   - An error if the query fails
func (r *actionRepositoryImpl) GetScheduledActionByID(ctx context.Context, actionID string) (db.ActionDBResponse, error) {
	var action db.ActionDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"id": actionID,
		},
	}

	if err := r.db.Select(action, "schedule_full_view", opts, "id", "*"); err != nil {
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
	var actionList []actions.Action
	for i := range newActions {
		actionList = append(actionList, newActions[i])
	}
	schedActions, cronActions := actions.SplitActionsByType(actionList)

	tx, err := r.db.BeginTxx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Writing Scheduled Actions
	if len(schedActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, InsertScheduledActionsQuery, schedActions); err != nil {
			// a.logger.Error("Failed to prepare InsertScheduledActionsQuery query", zap.Error(err))
			return fmt.Errorf("failed to insert scheduled actions: %w", err)
		}
	}

	// Writing Cron Actions
	if len(cronActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, InsertCronActionsQuery, cronActions); err != nil {
			// a.logger.Error("Failed to prepare InsertCronActionsQuery query", zap.Error(err))
			return fmt.Errorf("failed to insert cron actions: %w", err)
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// DeleteScheduledAction removes an actions.ScheduledAction action from the DB based on its ID
//
// Parameters:
//   - A string containing the action ID to be removed
//
// Returns:
//   - An error if the delete query fails
func (r *actionRepositoryImpl) DeleteScheduledAction(ctx context.Context, actionID string) error {
	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"id": actionID,
		},
	}

	if err := r.db.Delete("schedule", opts); err != nil {
		return err
	}
	return nil
}
