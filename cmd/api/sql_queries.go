package main

const (
	// SelectExpensesQuery returns every expense in the inventory ordered by instanceID
	SelectExpensesQuery = `
		SELECT * FROM expenses
		ORDER BY instance_id
	`

	// SelectLastExpensesQuery returns the last expense for every instance older
	// than 1 day. This is used for obtainning the list of instances that need
	// Billing information update because all the instances returned by this
	// query doesn't have expenses for the current day
	SelectLastExpensesQuery = `
		SELECT
				instances.id
		FROM
				instances
		LEFT JOIN (
				SELECT
						instance_id,
						MAX(date) AS last_expense_date
				FROM
						expenses
				GROUP BY
						instance_id
		) AS last_expenses
		ON
				instances.id = last_expenses.instance_id
		WHERE
				last_expenses.last_expense_date IS NULL
				OR last_expenses.last_expense_date < CURRENT_DATE - INTERVAL '1 day';
	`
	SelectLastExpensesQueryVOPT = `
		WITH ranked_expenses AS (
				SELECT
						instance_id,
						date,
						amount,
						ROW_NUMBER() OVER (PARTITION BY instance_id ORDER BY date DESC) AS rn
				FROM expenses
		)
		SELECT instance_id, date, amount FROM ranked_expenses JOIN instances ON instance_id = id WHERE rn = 1 AND (status != 'Terminated' AND status != 'Unknown') AND date < '$1';
	`

	// SelectExpensesByInstanceQuery returns expense in the inventory for a specific InstanceID
	SelectExpensesByInstanceQuery = `
		SELECT * FROM expenses
		WHERE instance_id = $1
		ORDER BY date
	`

	// InsertExpensesQuery inserts into a new expense for an instance
	InsertExpensesQuery = `
		INSERT INTO expenses (
			instance_id,
			date,
			amount
		) VALUES (
			:instance_id,
			:date,
			:amount
		) ON CONFLICT (instance_id, date) DO UPDATE SET
			amount = EXCLUDED.amount
	`

	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = `
		SELECT * FROM instances
		JOIN tags ON
			instances.id = tags.instance_id
		ORDER BY name
	`

	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = `
		SELECT * FROM instances
		JOIN tags ON
			instances.id = tags.instance_id
		WHERE id = $1
		ORDER BY name
	`

	// SelectClustersQuery returns every cluster in the inventory ordered by Name
	SelectClustersQuery = `
		SELECT * FROM clusters
		ORDER BY name
	`
	// InsertEvent insert a new audit event
	InsertEvent = `
    INSERT INTO audit_log(
        event_timestamp,
        triggered_by,
        action_name,
        resource_id,
        resource_type,
        result,
        description,
        severity
    ) VALUES (
        CURRENT_TIMESTAMP,
        :triggered_by,
        :action_name,
        :resource_id,
        :resource_type,
        :result,
        :description,
        :severity
    ) RETURNING id`
	// SelectClusterEvents returns audit log events related to a specific cluster.
	SelectClusterEvents = `
	SELECT 
		al.id, 
		al.event_timestamp, 
		al.triggered_by, 
		al.action_name,
		al.resource_id, 
		al.resource_type, 
		al.result, 
		al.description, 
		al.severity
	FROM audit_log al
	JOIN clusters c ON al.resource_id = c.id
	WHERE c.id = $1 AND al.resource_type = 'cluster'
	`

	// SelectClusterAccountNameQuery returns an cluster by its Name
	SelectClusterAccountNameQuery = `
		SELECT account_name FROM clusters
		WHERE id = $1
	`

	// SelectClusterRegionQuery returns an cluster by its Name
	SelectClusterRegionQuery = `
		SELECT region FROM clusters
		WHERE id = $1
	`

	// SelectClustersByIDuery returns an cluster by its Name
	SelectClustersByIDuery = `
		SELECT * FROM clusters
		WHERE id = $1
		ORDER BY name
	`

	// SelectClusterTags returns the cluster's tags
	SelectClusterTags = `
		SELECT DISTINCT ON (key) key,value,instance_id FROM instances
		JOIN tags ON
			id=instance_id
		WHERE cluster_id = $1
	`

	// SelectInstancesOnClusterQuery returns every instance belonging to a acluster
	SelectInstancesOnClusterQuery = `
		SELECT * FROM instances
		WHERE cluster_id = $1
		ORDER BY id
	`

	// SelectAccountsQuery returns every instance in the inventory ordered by Name
	SelectAccountsQuery = `
		SELECT * FROM accounts
		ORDER BY name
	`

	// SelectAccountsByNameQuery returns an instance by its Name
	SelectAccountsByNameQuery = `
		SELECT * FROM accounts
		WHERE name = $1
		ORDER BY name
	`

	// SelectClustersOnAccountQuery returns an cluster by its Name
	SelectClustersOnAccountQuery = `
		SELECT * FROM clusters
		WHERE account_name = $1
		ORDER BY name
	`

	// InsertInstancesQuery inserts into a new instance in its table
	InsertInstancesQuery = `
		INSERT INTO instances (
			id,
			name,
			provider,
			instance_type,
			availability_zone,
			status,
			cluster_id,
			last_scan_timestamp,
			creation_timestamp,
			age,
			daily_cost,
			total_cost
		) VALUES (
			:id,
			:name,
			:provider,
			:instance_type,
			:availability_zone,
			:status,
			:cluster_id,
			:last_scan_timestamp,
			:creation_timestamp,
			:age,
			:daily_cost,
			:total_cost
		) ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			provider = EXCLUDED.provider,
			instance_type = EXCLUDED.instance_type,
			availability_zone = EXCLUDED.availability_zone,
			status = EXCLUDED.status,
			cluster_id = EXCLUDED.cluster_id,
			last_scan_timestamp = EXCLUDED.last_scan_timestamp,
			creation_timestamp = EXCLUDED.creation_timestamp,
			age = EXCLUDED.age
	`

	// InsertClustersQuery inserts into a new instance in its table
	InsertClustersQuery = `
		INSERT INTO clusters (
			id,
			name,
			infra_id,
			provider,
			status,
			region,
			account_name,
			console_link,
			instance_count,
			last_scan_timestamp,
			creation_timestamp,
			age,
			owner,
			total_cost
		) VALUES (
			:id,
			:name,
			:infra_id,
			:provider,
			:status,
			:region,
			:account_name,
			:console_link,
			:instance_count,
			:last_scan_timestamp,
			:creation_timestamp,
			:age,
			:owner,
			:total_cost
		) ON CONFLICT (id) DO UPDATE SET
			provider = EXCLUDED.provider,
			status = EXCLUDED.status,
			region = EXCLUDED.region,
			console_link = EXCLUDED.console_link,
			instance_count = EXCLUDED.instance_count,
			last_scan_timestamp = EXCLUDED.last_scan_timestamp,
			creation_timestamp = EXCLUDED.creation_timestamp,
			age = EXCLUDED.age,
			owner = EXCLUDED.owner
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

	// InsertTagsQuery inserts into a new tag for an instance
	InsertTagsQuery = `
		INSERT INTO tags (
			key,
			value,
			instance_id
		) VALUES (
			:key,
			:value,
			:instance_id
		) ON CONFLICT (key, instance_id) DO UPDATE SET
			value = EXCLUDED.value
	`

	// DeleteInstanceQuery removes an instance by its ID
	DeleteInstanceQuery = `DELETE FROM instances WHERE id=$1`

	// DeleteClusterQuery removes an cluster by its name
	DeleteClusterQuery = `DELETE FROM clusters WHERE id=$1`

	// DeleteAccountQuery removes an account by its name
	DeleteAccountQuery = `DELETE FROM accounts WHERE name=$1`

	// DeleteTagsQuery removes a Tag by its key and instance reference
	DeleteTagsQuery = `DELETE FROM tags WHERE instance_id=$1`

	// Run 'check_terminated_instances()' SQL function
	UpdateTerminatedInstancesQuery = `SELECT check_terminated_instances()`
	// Run 'check_terminated_clusters()' SQL function
	UpdateTerminatedClustersQuery = `SELECT check_terminated_clusters()`

	// UpdateInstanceStatus updates the status of a  set of instances based on their clusterID
	UpdateStatusClusterByClusterIDQuery = `UPDATE clusters SET status=$1 WHERE id=$2`

	// UpdateInstanceStatus updates the status of a  set of instances based on their clusterID
	UpdateStatusInstancesByClusterIDQuery = `UPDATE instances SET status=$1 WHERE cluster_id=$2`

	// CheckStatusQuery checks if the requested status exists on the DB
	CheckStatusQuery = `SELECT EXISTS (SELECT 1 FROM status WHERE value=$1)`
)
