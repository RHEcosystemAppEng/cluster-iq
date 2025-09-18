package repositories

const (
	// Table for SELECT operations on Accounts
	SelectAccountsView = "accounts_full_view"

	// InsertAccountsQuery inserts into a new instance in its table
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
