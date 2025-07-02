package repositories

const (
	// SelectAccountsByNameQuery returns an instance by its Name
	SelectAccountsByNameQuery = `
		SELECT * FROM accounts
		WHERE name = $1
		ORDER BY name
	`
	// InsertAccountsQuery inserts into a new instance in its table
	InsertAccountsQuery = `
		INSERT INTO accounts (
			id,
			name,
			provider,
			total_cost,
			cluster_count,
			last_scan_timestamp
		) VALUES (
			:id,
			:name,
			:provider,
			:total_cost,
			:cluster_count,
			:last_scan_timestamp
		) ON CONFLICT (name) DO UPDATE SET
			id = EXCLUDED.id,
			provider = EXCLUDED.provider,
			cluster_count = EXCLUDED.cluster_count,
			last_scan_timestamp = EXCLUDED.last_scan_timestamp
	`
	// DeleteAccountQuery removes an account by its name
	DeleteAccountQuery = `DELETE FROM accounts WHERE name=$1`
	// SelectAccountsQuery returns every instance in the inventory ordered by Name
	SelectAccountsQuery = `SELECT * FROM accounts`
)
