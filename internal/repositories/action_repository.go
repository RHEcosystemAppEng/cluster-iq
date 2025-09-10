package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var _ ActionRepository = (*actionRepositoryImpl)(nil)

// ScheduleRecord is a helper struct to scan the result from the database
// before converting it to a specific action type.
type ScheduleRecord struct {
	ID          string         `db:"id"`
	Type        string         `db:"type"`
	Time        sql.NullTime   `db:"time"`
	CronExp     sql.NullString `db:"cron_exp"`
	Operation   string         `db:"operation"`
	Status      string         `db:"status"`
	Enabled     bool           `db:"enabled"`
	ClusterID   string         `db:"cluster_id"`
	Region      string         `db:"region"`
	AccountName string         `db:"account_name"`
	Instances   pq.StringArray `db:"instances"`
}

// toAction converts a ScheduleRecord to a concrete actions.Action implementation.
func (sr *ScheduleRecord) toAction() (actions.Action, error) {
	actionOp := actions.ActionOperation(sr.Operation)
	target := actions.ActionTarget{
		AccountName: sr.AccountName,
		Region:      sr.Region,
		ClusterID:   sr.ClusterID,
		Instances:   sr.Instances,
	}

	baseAction := actions.BaseAction{
		ID:        sr.ID,
		Operation: actionOp,
		Target:    target,
		Status:    sr.Status,
		Enabled:   sr.Enabled,
	}

	switch sr.Type {
	case string(actions.ScheduledActionType):
		if !sr.Time.Valid {
			return nil, fmt.Errorf("scheduled action %s has invalid time", sr.ID)
		}
		return &actions.ScheduledAction{
			BaseAction: baseAction,
			When:       sr.Time.Time,
			Type:       sr.Type,
		}, nil
	case string(actions.CronActionType):
		if !sr.CronExp.Valid {
			return nil, fmt.Errorf("cron action %s has invalid cron expression", sr.ID)
		}
		return &actions.CronAction{
			BaseAction: baseAction,
			Expression: sr.CronExp.String,
			Type:       sr.Type,
		}, nil
	default:
		return nil, fmt.Errorf("unknown action type: %s", sr.Type)
	}
}

// ActionRepository defines the interface for data access operations for actions.
type ActionRepository interface {
	ListScheduledActions(ctx context.Context, opts ListOptions) ([]actions.Action, int, error)
	EnableScheduledAction(ctx context.Context, actionID string) error
	DisableScheduledAction(ctx context.Context, actionID string) error
	GetScheduledActionByID(ctx context.Context, actionID string) (actions.Action, error)
	DeleteScheduledAction(ctx context.Context, actionID string) error
	Create(ctx context.Context, newActions []actions.Action) error
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
	var scheduleRecords []ScheduleRecord
	baseQuery := SelectScheduledActionsQuery
	countQuery := "SELECT COUNT(*) FROM schedule"

	whereClauses, namedArgs := buildActionWhereClauses(opts.Filters)

	total, err := listQueryHelper(ctx, r.db, &scheduleRecords, baseQuery, countQuery, opts, whereClauses, namedArgs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list scheduled actions: %w", err)
	}

	resultActions := make([]actions.Action, len(scheduleRecords))
	for i, record := range scheduleRecords {
		action, err := (&record).toAction()
		if err != nil {
			// Potentially log the error and continue, or fail fast
			return nil, 0, fmt.Errorf("failed to convert db action to action: %w", err)
		}
		resultActions[i] = action
	}

	return resultActions, total, nil
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
	defer func() { _ = tx.Rollback() }()

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
	defer func() { _ = tx.Rollback() }()

	// Writing Scheduled Actions
	if _, err := tx.ExecContext(ctx, DisableActionQuery, actionID); err != nil {
		// a.logger.Error("Failed to prepare DisableScheduledAction query", zap.Error(err))
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
	var scheduleRecord ScheduleRecord
	if err := r.db.GetContext(ctx, &scheduleRecord, SelectScheduledActionsByIDQuery, actionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get scheduled action by id %s: %w", actionID, err)
	}

	return (&scheduleRecord).toAction()
}

// Create creates a batch of scheduled actions in the database.
//
// Parameters:
//   - An array of actions.Action to write on the DB
//
// Returns:
//   - An error if the insert fails
func (r *actionRepositoryImpl) Create(ctx context.Context, newActions []actions.Action) error {
	// Begin transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	schedActions, cronActions := actions.SplitActionsByType(newActions)

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
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

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
		case "status": //nolint:goconst
			clauses = append(clauses, "status = :status") //nolint:goconst
			args["status"] = value
		// TODO. will it work with boolean?
		case "enabled":
			clauses = append(clauses, "enabled = :enabled")
			args["enabled"] = value
		}
	}

	return clauses, args
}
