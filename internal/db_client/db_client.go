// TODO: Placeholder for the SQL client to fix linter issues
// TODO: Add actual implementation in next PR
package dbclient

import (
	"context"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// DBClient defines the SQL interface for the API to interact with the database.
// It manages database connections and provides methods for interacting with various entities like instances, clusters, accounts, and expenses.
type DBClient struct {
	// db is the database connection object.
	db *sqlx.DB
	// logger is used for logging database operations and errors.
	logger *zap.Logger
	// Transactions context
	ctx context.Context
}

// NewDBClient initializes a new DBClient with the given database URL and logger.
//
// Parameters:
// - dbURL: The connection string for the PostgreSQL database.
// - logger: Logger instance for logging.
//
// Returns:
// - A pointer to an DBClient instance.
// - An error if the database connection fails.
func NewDBClient(dbURL string, logger *zap.Logger, ctx context.Context) (*DBClient, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return &DBClient{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}, nil
}

func (d *DBClient) Close() error {
	return d.db.Close()
}

func (d *DBClient) BeginTxx() (*sqlx.Tx, error) {
	return d.db.BeginTxx(d.ctx, nil)
}

func (d *DBClient) Ping() error {
	return d.db.Ping()
}

func (d *DBClient) Get(dest interface{}, table string, opts models.ListOptions, columns ...string) error {
	builder := d.NewSelectBuilder(columns...).From(table)

	// Preparing positional arguments
	i := 1
	for k, v := range opts.Filters {
		cond := fmt.Sprintf("%s = $%d", k, i)
		builder = builder.Where(cond, v)
		i++
	}

	query, args, err := builder.Build()
	d.logger.Debug("SELECT QUERY", zap.String("query", query))
	if err != nil {
		d.logger.Error("Error building SELECT query", zap.String("query", query), zap.Reflect("args", args), zap.Error(err))
		return err
	}
	return d.db.Get(dest, query, args...)
}

func (d *DBClient) Select(dest interface{}, table string, opts models.ListOptions, orderColumn string, columns ...string) error {
	builder := d.NewSelectBuilder(columns...).From(table)

	// Preparing positional arguments
	i := 1
	for k, v := range opts.Filters {
		cond := fmt.Sprintf("%s = $%d", k, i)
		builder = builder.Where(cond, v)
		i++
	}

	// Apply pagination
	if opts.PageSize > 0 {
		builder.Paginate(opts.PageSize, opts.Offset)
		builder.orderBy = orderColumn
	}

	query, args, err := builder.Build()
	d.logger.Debug("SELECT QUERY", zap.String("query", query))
	if err != nil {
		d.logger.Error("Error building SELECT query", zap.String("query", query), zap.Reflect("args", args), zap.Error(err))
		return err
	}
	return d.db.Select(dest, query, args...)
}

func (d *DBClient) Insert(query string, data interface{}) error {
	builder := d.NewInsertBuilder().Query(query).Data(data)

	tx, err := d.db.BeginTxx(d.ctx, nil)
	if err != nil {
		return err
	}

	// Rollback defer func
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				d.logger.Error("failed to Rollback INSERT")
			}
		}
	}()

	if _, err := tx.NamedExecContext(d.ctx, builder.query, builder.data); err != nil {
		return fmt.Errorf("named-exec INSERT error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit INSERT error: %w", err)
	}

	return nil
}
func (d *DBClient) Update(query string, data interface{}) error {
	builder := d.NewUpdateBuilder().Query(query).Data(data)

	tx, err := d.db.BeginTxx(d.ctx, nil)
	if err != nil {
		return err
	}

	// Rollback defer func
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				d.logger.Error("failed to Rollback UPDATE")
			}
		}
	}()

	if _, err := tx.ExecContext(d.ctx, builder.query, builder.data); err != nil {
		return fmt.Errorf("exec UPDATE error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit UPDATE error: %w", err)
	}

	return nil
}

func (d *DBClient) NamedUpdate(query string, data interface{}) error {
	builder := d.NewUpdateBuilder().Query(query).Data(data)

	tx, err := d.db.BeginTxx(d.ctx, nil)
	if err != nil {
		return err
	}

	// Rollback defer func
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				d.logger.Error("failed to Rollback UPDATE")
			}
		}
	}()

	if _, err := tx.NamedExecContext(d.ctx, builder.query, builder.data); err != nil {
		return fmt.Errorf("named-exec UPDATE error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit UPDATE error: %w", err)
	}

	return nil
}

// Delete executes a DELETE with a safe transaction pattern.Delete
func (d *DBClient) Delete(table string, opts models.ListOptions) error {
	builder := d.NewDeleteBuilder().From(table)

	// Processing "WHERE" conditions
	i := 1
	for k, v := range opts.Filters {
		cond := fmt.Sprintf("%s = $%d", k, i)
		builder = builder.Where(cond, v)
		i++
	}

	// Building query
	query, args, err := builder.Build()
	if err != nil {
		return fmt.Errorf("error building DELETE query: '%s'", query)
	}

	tx, err := d.db.BeginTxx(d.ctx, nil)
	if err != nil {
		return err
	}

	// Rollback defer func
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				d.logger.Error("failed to Rollback DELETE")
			}
		}
	}()

	if _, err := tx.ExecContext(d.ctx, query, args...); err != nil {
		return fmt.Errorf("exec DELETE error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit DELETE error: %w", err)
	}

	return nil
}
