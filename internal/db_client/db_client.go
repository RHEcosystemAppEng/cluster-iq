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
func NewDBClient(dbURL string, logger *zap.Logger) (*DBClient, error) {
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return &DBClient{
		db:     db,
		logger: logger,
	}, nil
}

func (d *DBClient) Close() error {
	return d.db.Close()
}

func (d *DBClient) NewTx(ctx context.Context) (*sqlx.Tx, error) {
	return d.db.BeginTxx(ctx, nil)
}

func (d *DBClient) Ping() error {
	return d.db.Ping()
}

func (d *DBClient) ExecFunc(ctx context.Context, query string) error {
	var result string
	if err := d.db.QueryRowxContext(ctx, query).Scan(&result); err != nil {
		return err
	}

	return nil
}

func (d *DBClient) GetWithContext(ctx context.Context, dest interface{}, table string, opts models.ListOptions, columns ...string) error {
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
	return d.db.GetContext(ctx, dest, query, args...)
}

func (d *DBClient) Get(dest interface{}, table string, opts models.ListOptions, columns ...string) error {
	return d.GetWithContext(context.TODO(), dest, table, opts, columns...)
}

// QueryRowContext executes a raw SQL query and scans the result into dest.
func (d *DBClient) QueryRowContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.db.GetContext(ctx, dest, query, args...)
}

func (d *DBClient) SelectWithContext(ctx context.Context, dest interface{}, table string, opts models.ListOptions, orderColumn string, columns ...string) error {
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
	return d.db.SelectContext(ctx, dest, query, args...)
}

func (d *DBClient) Select(dest interface{}, table string, opts models.ListOptions, orderColumn string, columns ...string) error {
	return d.SelectWithContext(context.TODO(), dest, table, opts, orderColumn, columns...)
}

func (d *DBClient) InsertWithReturnWithContext(ctx context.Context, query string, data interface{}) (int64, error) {
	builder := d.NewInsertBuilder().Query(query).Data(data)

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return -1, err
	}

	// Rollback defer func
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				d.logger.Error("failed to Rollback INSERT")
			}
		}
	}()

	var returnedValue int64
	rows, err := tx.NamedQuery(builder.query, builder.data)
	if err != nil {
		return -1, fmt.Errorf("named-exec INSERT error: %w", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			d.logger.Error("failed to close rows after insert")
		}
	}()

	if rows.Next() {
		if err := rows.Scan(&returnedValue); err != nil {
			return -1, fmt.Errorf("scan INSERT return value error %w", err)
		}
	} else {
		return -1, fmt.Errorf("sql INSERT did not return any value")
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("commit INSERT error: %w", err)
	}

	return returnedValue, nil
}

func (d *DBClient) InsertWithContext(ctx context.Context, query string, data interface{}) error {
	builder := d.NewInsertBuilder().Query(query).Data(data)

	tx, err := d.db.BeginTxx(ctx, nil)
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

	if _, err := tx.NamedExecContext(ctx, builder.query, builder.data); err != nil {
		return fmt.Errorf("named-exec INSERT error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit INSERT error: %w", err)
	}

	return nil
}

func (d *DBClient) Insert(query string, data interface{}) error {
	return d.InsertWithContext(context.TODO(), query, data)
}

func (d *DBClient) UpdateWithContext(ctx context.Context, query string, data interface{}) error {
	builder := d.NewUpdateBuilder().Query(query).Data(data)

	tx, err := d.db.BeginTxx(ctx, nil)
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

	if _, err := tx.ExecContext(ctx, builder.query, builder.data); err != nil {
		return fmt.Errorf("exec UPDATE error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit UPDATE error: %w", err)
	}

	return nil
}

func (d *DBClient) Update(query string, data interface{}) error {
	return d.UpdateWithContext(context.TODO(), query, data)
}

func (d *DBClient) NamedUpdateWithContext(ctx context.Context, query string, data interface{}) error {
	builder := d.NewUpdateBuilder().Query(query).Data(data)

	tx, err := d.db.BeginTxx(ctx, nil)
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

	if _, err := tx.NamedExecContext(ctx, builder.query, builder.data); err != nil {
		return fmt.Errorf("named-exec UPDATE error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit UPDATE error: %w", err)
	}

	return nil
}

func (d *DBClient) NamedUpdate(query string, data interface{}) error {
	return d.NamedUpdateWithContext(context.TODO(), query, data)
}

// Delete executes a DELETE with a safe transaction pattern.Delete
func (d *DBClient) DeleteWithContext(ctx context.Context, table string, opts models.ListOptions) error {
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

	tx, err := d.db.BeginTxx(ctx, nil)
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

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("exec DELETE error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit DELETE error: %w", err)
	}

	return nil
}

func (d *DBClient) Delete(table string, opts models.ListOptions) error {
	return d.DeleteWithContext(context.TODO(), table, opts)
}
