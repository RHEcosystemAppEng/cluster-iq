package sqlclient

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	dbmodel "github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"go.uber.org/zap"
)

const (
	// SelectAccountsQuery returns every instance in the inventory ordered by Name
	SelectAccountsQuery = `
		SELECT * FROM accounts_full_view
		ORDER BY account_name
	`

	// SelectAccountsByIDQuery returns an instance by its Name
	SelectAccountsByIDQuery = `
		SELECT * FROM accounts_full_view
		WHERE account_id = $1
		ORDER BY account_name
	`

	// InsertAccountsQuery inserts into a new instance in its table
	InsertAccountsQuery = `
		INSERT INTO accounts (
			account_id,
			account_name,
			provider,
			last_scan_ts
		) VALUES (
			:account_id,
			:account_name,
			:provider,
			:last_scan_ts
		) ON CONFLICT (account_id) DO UPDATE SET
			account_name = EXCLUDED.account_name,
			provider = EXCLUDED.provider,
			last_scan_ts = EXCLUDED.last_scan_ts
	`

	// SelectClustersOnAccountQuery returns an cluster by its Name
	SelectClustersOnAccountQuery = `
		SELECT * FROM clusters_full_view
		WHERE account_id = $1
		ORDER BY cluster_name
	`

	// DeleteAccountQuery removes an account by its name
	DeleteAccountQuery = `DELETE FROM accounts WHERE account_id = $1`
)

// GetAccounts retrieves all accounts from the database.
//
// Returns:
// - A slice of model.AccountDBResponse objects.
// - An error if the query fails.
func (a SQLClient) GetAccounts() ([]dbmodel.AccountDBResponse, error) {
	var accounts []dbmodel.AccountDBResponse
	if err := a.db.Select(&accounts, SelectAccountsQuery); err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountByID retrieves an account by its ID from the database.
//
// Parameters:
// - accountID: The name of the account to retrieve.
//
// Returns:
// - A pointer of dbmodel.AccountDBResponse.
// - An error if the query fails.
func (a SQLClient) GetAccountByID(accountID string) (*dbmodel.AccountDBResponse, error) {
	var account dbmodel.AccountDBResponse
	if err := a.db.Get(&account, SelectAccountsByIDQuery, accountID); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetClustersOnAccount retrieves all clusters associated with a specific account.
//
// Parameters:
// - accountID: The name of the account whose clusters will be retrieved.
//
// Returns:
// - A slice of inventory.Cluster objects.
// - An error if the query fails.
func (a SQLClient) GetClustersOnAccount(accountID string) ([]dbmodel.ClusterDBResponse, error) {
	var clusters []dbmodel.ClusterDBResponse
	if err := a.db.Select(&clusters, SelectClustersOnAccountQuery, accountID); err != nil {
		return nil, err
	}
	return clusters, nil
}

// WriteAccounts inserts multiple accounts into the database in a transaction.
//
// Parameters:
// - accounts: A slice of inventory.Account objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (a SQLClient) WriteAccounts(accounts []inventory.Account) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback WriteAccounts transaction", zap.Error(rbErr))
			}
		}
	}()

	if _, err = tx.NamedExec(InsertAccountsQuery, accounts); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// DeleteAccount deletes an account from the database by its name.
//
// Parameters:
// - accountID: The name of the account to delete.
//
// Returns:
// - An error if the transaction fails.
func (a SQLClient) DeleteAccount(accountID string) error {
	tx, err := a.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				a.logger.Error("Failed to rollback Delete Account transaction", zap.Error(rbErr))
			}
		}
	}()

	tx.MustExec(DeleteAccountQuery, accountID)
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
