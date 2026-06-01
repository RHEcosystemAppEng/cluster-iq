package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
)

const (
	ActionRunsTable = "action_runs"

	InsertActionRunQuery = `
		INSERT INTO action_runs (schedule_id)
		VALUES ($1)
		RETURNING id
	`

	UpdateActionRunQuery = `
		UPDATE action_runs
		SET
			status = $1,
			finished_at = NOW(),
			error_msg = $2
		WHERE id = $3
	`
)

var _ ActionRunRepository = (*actionRunRepositoryImpl)(nil)

// ActionRunRepository defines the interface for action run data access operations.
type ActionRunRepository interface {
	List(ctx context.Context, opts models.ListOptions) ([]db.ActionRunDBResponse, int, error)
	GetByID(ctx context.Context, runID string) (db.ActionRunDBResponse, error)
	Create(ctx context.Context, scheduleID string) (int64, error)
	Update(ctx context.Context, runID string, status string, errorMsg string) error
}

type actionRunRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewActionRunRepository(db *dbclient.DBClient) ActionRunRepository {
	return &actionRunRepositoryImpl{db: db}
}

func (r *actionRunRepositoryImpl) List(ctx context.Context, opts models.ListOptions) ([]db.ActionRunDBResponse, int, error) {
	runs := []db.ActionRunDBResponse{}

	if err := r.db.SelectWithContext(ctx, &runs, ActionRunsTable, opts, "id", "*"); err != nil {
		return runs, 0, fmt.Errorf("failed to list action runs: %w", err)
	}

	return runs, len(runs), nil
}

func (r *actionRunRepositoryImpl) GetByID(ctx context.Context, runID string) (db.ActionRunDBResponse, error) {
	var run db.ActionRunDBResponse

	opts := models.ListOptions{
		Filters: map[string]interface{}{
			"id": runID,
		},
	}

	if err := r.db.GetWithContext(ctx, &run, ActionRunsTable, opts, "id", "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return run, ErrNotFound
		}
		return run, err
	}

	return run, nil
}

func (r *actionRunRepositoryImpl) Create(ctx context.Context, scheduleID string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, &id, InsertActionRunQuery, scheduleID)
	if err != nil {
		return -1, fmt.Errorf("failed to create action run: %w", err)
	}
	return id, nil
}

func (r *actionRunRepositoryImpl) Update(ctx context.Context, runID string, status string, errorMsg string) error {
	tx, err := r.db.NewTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, UpdateActionRunQuery, status, errorMsg, runID); err != nil {
		return fmt.Errorf("failed to update action run: %w", err)
	}

	return tx.Commit()
}
