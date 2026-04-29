package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
)

const (
	// DB Table for accounts
	AccountsTable = "accounts"
	// View for SELECT operations on Accounts
	SelectAccountsView = "accounts_full_view"
	// Materialized view for SELECT operations on Accounts
	SelectAccountsMView = "m_accounts_full_view"
	// View for SELECT operations for Instances pending on Expense Update
	SelectInstancesPendingExpenseUpdateView = "instances_pending_expense_update"
	// InsertAccountsQuery to insert or update new accounts
	InsertAccountsQuery = `
		INSERT INTO accounts (
			account_id,
			account_name,
			provider,
			last_scan_ts,
			created_at
		) VALUES (
			:account_id,
			:account_name,
			:provider,
			:last_scan_ts,
			:created_at
		) ON CONFLICT (account_id, provider) DO UPDATE SET
			last_scan_ts = EXCLUDED.last_scan_ts
	`
)

var _ AccountRepository = (*accountRepositoryImpl)(nil)

// AccountRepository defines the interface for data access operations for accounts.
type AccountRepository interface {
	ListAccounts(ctx context.Context, opts models.ListOptions) ([]db.AccountDBResponse, int, error)
	CountAccounts(ctx context.Context, opts models.ListOptions) (int, error)
	GetAccountByID(ctx context.Context, accountID string) (db.AccountDBResponse, error)
	GetAccountClustersByID(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error)
	GetExpenseUpdateInstances(ctx context.Context, accountID string) ([]db.InstanceDBResponse, error)
	GetScannerTimestamp(ctx context.Context) (time.Time, error)
	CreateAccount(ctx context.Context, accounts []inventory.Account) error
	UpdateAccount(ctx context.Context, accountID string, patch dto.AccountPatchRequest) error
	DeleteAccount(ctx context.Context, accountID string) error
}

type accountRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewAccountRepository(db *dbclient.DBClient) AccountRepository {
	return &accountRepositoryImpl{db: db}
}

// ListAccounts retrieves all accounts from the database.
//
// Returns:
// - A slice of inventory.Account objects.
// - An error if the query fails.
func (r *accountRepositoryImpl) ListAccounts(ctx context.Context, opts models.ListOptions) ([]db.AccountDBResponse, int, error) {
	accounts := []db.AccountDBResponse{}

	if err := r.db.SelectWithContext(ctx, &accounts, SelectAccountsMView, opts, "account_id", "*"); err != nil {
		return accounts, 0, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, len(accounts), nil
}

// CountAccounts performs counting operations.
//
// Returns:
// - The number of counted accounts.
// - An error if the query fails.
func (r *accountRepositoryImpl) CountAccounts(ctx context.Context, opts models.ListOptions) (int, error) {
	var count int

	if err := r.db.GetWithContext(ctx, &count, SelectAccountsMView, opts, "count(*)"); err != nil {
		return count, fmt.Errorf("failed to list accounts: %w", err)
	}

	return count, nil
}

// GetAccountByID retrieves an account by its name from the database.
//
// Parameters:
// - accountID: The name of the account to retrieve.
//
// Returns:
// - A slice of inventory.Account objects (usually containing one element).
// - An error if the query fails.
func (r *accountRepositoryImpl) GetAccountByID(ctx context.Context, accountID string) (db.AccountDBResponse, error) {
	var account db.AccountDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"account_id": accountID,
		},
	}

	if err := r.db.GetWithContext(ctx, &account, SelectAccountsView, opts, "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account, ErrNotFound
		}
		return account, err
	}
	return account, nil
}

// GetAccountClustersByID retrieves an account by its name from the database.
//
// Parameters:
// - accountID: The name of the account to retrieve.
//
// Returns:
// - A slice of inventory.Account objects (usually containing one element).
// - An error if the query fails.
func (r *accountRepositoryImpl) GetAccountClustersByID(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error) {
	if _, err := r.GetAccountByID(ctx, accountID); err != nil {
		return nil, err
	}

	clusters := []db.ClusterDBResponse{}

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"account_id": accountID,
		},
	}

	if err := r.db.SelectWithContext(ctx, &clusters, SelectClustersFullMView, opts, "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return clusters, nil
		}
		return clusters, err
	}

	if len(clusters) == 0 {
		return clusters, ErrNoClustersInAccount
	}

	return clusters, nil
}

// GetExpenseUpdateInstances retrieves instances with outdated billing information.
//
// Parameters:
//
// Returns:
// - A slice of inventory.Instance objects.
// - An error if the query fails.
func (r *accountRepositoryImpl) GetExpenseUpdateInstances(ctx context.Context, accountID string) ([]db.InstanceDBResponse, error) {
	instances := []db.InstanceDBResponse{}

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"account_id": accountID,
		},
	}

	if err := r.db.SelectWithContext(ctx, &instances, SelectInstancesPendingExpenseUpdateView, opts, "instance_id", "instance_id"); err != nil {
		return instances, fmt.Errorf("failed to list instances pending of expense update: %w", err)
	}

	return instances, nil
}

// Create inserts multiple accounts into the database in a transaction.
//
// Parameters:
// - accounts: A slice of inventory.Account objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (r *accountRepositoryImpl) CreateAccount(ctx context.Context, accounts []inventory.Account) error {
	if err := r.db.InsertWithContext(ctx, InsertAccountsQuery, accounts); err != nil {
		return err
	}

	return nil
}

// refreshAccountsMView refreshes the accounts materialized view in a separate transaction.
func (r *accountRepositoryImpl) refreshAccountsMView(ctx context.Context) error {
	tx, txErr := r.db.NewTx(ctx)
	if txErr != nil {
		return fmt.Errorf("failed to create transaction for refresh: %w", txErr)
	}

	var err error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.ExecContext(ctx, "REFRESH MATERIALIZED VIEW m_accounts_full_view"); err != nil {
		return fmt.Errorf("refresh materialized view error: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit refresh error: %w", err)
	}

	return nil
}

// UpdateAccount updates mutable fields of an existing account in the database.
// Only non-nil fields in the patch request will be updated.
func (r *accountRepositoryImpl) UpdateAccount(ctx context.Context, accountID string, patch dto.AccountPatchRequest) (err error) {
	// Build dynamic UPDATE query with positional parameters
	query := "UPDATE accounts SET "
	args := make([]interface{}, 0)
	argCount := 1

	if patch.AccountName != nil {
		query += fmt.Sprintf("account_name = $%d", argCount)
		args = append(args, *patch.AccountName)
		argCount++
	}

	// If no fields to update, return early
	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE account_id = $%d", argCount)
	args = append(args, accountID)

	// Execute update in a transaction
	tx, txErr := r.db.NewTx(ctx)
	if txErr != nil {
		return txErr
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, execErr := tx.ExecContext(ctx, query, args...); execErr != nil {
		err = fmt.Errorf("exec UPDATE error: %w", execErr)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit UPDATE error: %w", err)
	}

	// Refresh materialized view after transaction commits
	return r.refreshAccountsMView(ctx)
}

// DeleteAccount deletes an account from the database by its ID.
//
// Parameters:
// - accountID: The ID of the account to delete.
//
// Returns:
// - An error if the transaction fails.
func (r *accountRepositoryImpl) DeleteAccount(ctx context.Context, accountID string) error {
	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"account_id": accountID,
		},
	}

	if err := r.db.DeleteWithContext(ctx, "accounts", opts); err != nil {
		return err
	}
	return nil
}

// GetScannerTimestamp retrieves the latest scan timestamp from all accounts.
//
// Returns:
// - The latest scan timestamp.
// - An error if the query fails.
func (r *accountRepositoryImpl) GetScannerTimestamp(ctx context.Context) (time.Time, error) {
	var timestamp time.Time

	if err := r.db.QueryRowContext(ctx, &timestamp, "SELECT MAX(last_scan_ts) FROM accounts"); err != nil {
		return timestamp, err
	}
	return timestamp, nil
}
