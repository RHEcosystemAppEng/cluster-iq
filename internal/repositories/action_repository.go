package repositories

import (
	"context"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/jmoiron/sqlx"
)

var _ ActionRepository = (*actionRepositoryImpl)(nil)

// ActionRepository defines the interface for data access operations for actions.
type ActionRepository interface {
	ListScheduledActions(ctx context.Context, opts ListOptions) ([]actions.Action, int, error)
	EnableScheduledAction(ctx context.Context, actionID string) error
	DisableScheduledAction(ctx context.Context, actionID string) error
	GetScheduledActionByID(ctx context.Context, actionID string) (actions.Action, error)
	WriteScheduledActions(ctx context.Context, newActions []actions.Action) error
	PatchScheduledAction(ctx context.Context, newActions []actions.Action) error
	PatchScheduledActionStatus(ctx context.Context, actionID string, status string) error
	DeleteScheduledAction(ctx context.Context, actionID string) error
}

type actionRepositoryImpl struct {
	db *sqlx.DB
}

func NewActionRepository(db *sqlx.DB) ActionRepository {
	return &actionRepositoryImpl{db: db}
}

// ListScheduledActions runs the db select query for retrieving the scheduled actions on the DB
//
// Parameters:
//
// Returns:
//   - An array of actions.Action with the scheduled actions declared on the DB
//   - An error if the query fails
func (r *actionRepositoryImpl) ListScheduledActions(ctx context.Context, opts ListOptions) ([]actions.Action, int, error) {
	var scheduledActions []actions.Action
	// Assuming a generic query and a count query exist for actions
	baseQuery := SelectScheduledActionsQuery
	countQuery := "SELECT COUNT(*) FROM schedule"

	whereClauses, namedArgs := buildActionWhereClauses(opts.Filters)

	total, err := listQueryHelper(ctx, r.db, &scheduledActions, baseQuery, countQuery, opts, whereClauses, namedArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list scheduled actions: %w", err)
	}

	return scheduledActions, total, nil
}

// EnableScheduledAction enables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) EnableScheduledAction(ctx context.Context, actionID string) error {
	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Writing Scheduled Actions
	if _, err := tx.ExecContext(ctx, EnableActionQuery, actionID); err != nil {
		// a.logger.Error("Failed to prepare EnableScheduledAction query", zap.Error(err))
		return fmt.Errorf("failed to enable action %s: %w", actionID, err)
	}

	// Commit the transaction
	return tx.Commit()
}

// DisableScheduledAction Disables an Action by its ID
//
// Parameters:
//   - Action ID
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) DisableScheduledAction(ctx context.Context, actionID string) error {
	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Writing Scheduled Actions
	if _, err := tx.ExecContext(ctx, DisableActionQuery, actionID); err != nil {
		//a.logger.Error("Failed to prepare DisableScheduledAction query", zap.Error(err))
		return fmt.Errorf("failed to disable action %s: %w", actionID, err)
	}

	// Commit the transaction
	return tx.Commit()
}

// GetScheduledActionByID runs the db select query for retrieving a specific scheduled action by its ID
//
// Parameters:
//
// Returns:
//   - An actions.Action object
//   - An error if the query fails
func (r *actionRepositoryImpl) GetScheduledActionByID(ctx context.Context, actionID string) (actions.Action, error) {
	// Getting results from DB
	var action actions.Action
	if err := r.db.GetContext(ctx, &action, SelectScheduledActionsByIDQuery, actionID); err != nil {
		//a.logger.Error("Failed to prepare SelectScheduledActions query", zap.Error(err))
		return nil, fmt.Errorf("failed to get scheduled action by id %s: %w", actionID, err)
	}

	return action, nil
}

// WriteScheduledActions receives an array of actions.Action and writes them on the DB
//
// Parameters:
//   - An array of actions.Action to write on the DB
//
// Returns:
//   - An error if the insert fails
func (r *actionRepositoryImpl) WriteScheduledActions(ctx context.Context, newActions []actions.Action) error {
	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	schedActions, cronActions := actions.SplitActionsByType(newActions)

	// Writing Scheduled Actions
	if len(schedActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, InsertScheduledActionsQuery, schedActions); err != nil {
			//a.logger.Error("Failed to prepare InsertScheduledActionsQuery query", zap.Error(err))
			return fmt.Errorf("failed to insert scheduled actions: %w", err)
		}
	}

	// Writing Cron Actions
	if len(cronActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, InsertCronActionsQuery, cronActions); err != nil {
			//a.logger.Error("Failed to prepare InsertCronActionsQuery query", zap.Error(err))
			return fmt.Errorf("failed to insert cron actions: %w", err)
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// PatchScheduledAction updates Action by its ID
//
// Parameters:
//   - Action list
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) PatchScheduledAction(ctx context.Context, newActions []actions.Action) error {
	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	schedActions, cronActions := actions.SplitActionsByType(newActions)

	// Writing Scheduled Actions
	if len(schedActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, PatchScheduledActionsQuery, schedActions); err != nil {
			//a.logger.Error("Failed to prepare PatchScheduledAction query", zap.Error(err))
			return fmt.Errorf("failed to patch scheduled actions: %w", err)
		}
	}

	// Writing Cron Actions
	if len(cronActions) > 0 {
		if _, err := tx.NamedExecContext(ctx, PatchCronActionsQuery, cronActions); err != nil {
			//a.logger.Error("Failed to prepare PatchCronActionsQuery query", zap.Error(err))
			return fmt.Errorf("failed to patch cron actions: %w", err)
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// PatchScheduledActionStatus updates Action status by its ID
//
// Parameters:
//   - Action list
//
// Returns:
//   - An error if the query fails
func (r *actionRepositoryImpl) PatchScheduledActionStatus(ctx context.Context, actionID string, status string) error {
	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	enabled := status == "Pending"

	if _, err := tx.ExecContext(ctx, PatchActionStatusQuery, actionID, status, enabled); err != nil {
		//a.logger.Error("Failed to prepare PatchScheduledActionStatus query", zap.Error(err))
		return fmt.Errorf("failed to patch action status for %s: %w", actionID, err)
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
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Deleting
	if _, err := tx.ExecContext(ctx, DeleteScheduledActionsQuery, actionID); err != nil {
		return fmt.Errorf("failed to delete scheduled action %s: %w", actionID, err)
	}

	// Commit the transaction
	return tx.Commit()
}

func buildActionWhereClauses(filters map[string]interface{}) ([]string, map[string]interface{}) {
	clauses := make([]string, 0, len(filters))
	args := make(map[string]interface{})

	for key, value := range filters {
		switch key {
		case "type":
			clauses = append(clauses, "type = :type")
			args["type"] = value
		case "operation":
			clauses = append(clauses, "operation = :operation")
			args["operation"] = value
		case "status":
			clauses = append(clauses, "status = :status")
			args["status"] = value
		// TODO. will it work with boolean?
		case "enabled":
			clauses = append(clauses, "enabled = :enabled")
			args["enabled"] = value
		}
	}

	return clauses, args
}
