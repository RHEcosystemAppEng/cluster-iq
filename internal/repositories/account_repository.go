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
	GetAccountByID(ctx context.Context, accountID string) (db.AccountDBResponse, error)
	GetAccountClustersByID(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error)
	GetExpenseUpdateInstances(ctx context.Context, accountID string) ([]db.InstanceDBResponse, error)
	GetScannerTimestamp(ctx context.Context) (time.Time, error)
	CreateAccount(ctx context.Context, accounts []inventory.Account) error
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
	var accounts []db.AccountDBResponse

	if err := r.db.SelectWithContext(ctx, &accounts, SelectAccountsMView, opts, "account_id", "*"); err != nil {
		return accounts, 0, fmt.Errorf("failed to list accounts: %w", err)
	}

	return accounts, len(accounts), nil
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

	if err := r.db.GetWithContext(ctx, &account, SelectAccountsMView, opts, "*"); err != nil {
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
	var clusters []db.ClusterDBResponse

	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"account_id": accountID,
		},
	}

	if err := r.db.SelectWithContext(ctx, &clusters, SelectClustersFullMView, opts, "*"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return clusters, ErrNotFound
		}
		return clusters, err
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
	var instances []db.InstanceDBResponse

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
