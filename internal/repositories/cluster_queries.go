package repositories

const (
	// Table for SELECT operations on Accounts
	SelectClustersFullView = "clusters_full_view"

	// UpdateStatusClusterByClusterIDQuery updates the status of a  set of instances based on their clusterID
	UpdateStatusClusterByClusterIDQuery = `UPDATE clusters SET status=$1 WHERE id=$2`
	// SelectClustersQuery returns every cluster in the inventory ordered by Name
	SelectClustersQuery = `SELECT * FROM clusters`
	// SelectClustersOverview returns the number of clusters grouped by status
	SelectClustersOverview = `
		SELECT
			COUNT(CASE WHEN status = 'Running' THEN 1 END) AS running,
			COUNT(CASE WHEN status = 'Stopped' THEN 1 END) AS stopped,
			COUNT(CASE WHEN status = 'Terminated' THEN 1 END) AS archived
		FROM clusters;
	`
	// SelectClusterAccountNameQuery returns a cluster by its Name
	SelectClusterAccountNameQuery = `
		SELECT account_name FROM clusters
		WHERE id = $1
	`
	// SelectClusterRegionQuery returns a cluster by its Name
	SelectClusterRegionQuery = `
		SELECT region FROM clusters
		WHERE id = $1
	`
	// SelectClustersByIDQuery returns a cluster by its Name
	SelectClustersByIDQuery = `
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
	// SelectClustersOnAccountQuery returns a cluster by its Name
	SelectClustersOnAccountQuery = `
		SELECT * FROM clusters
		WHERE account_name = $1
		ORDER BY name
	`
	// SelectInstancesOnClusterQuery returns every instance belonging to a cluster
	SelectInstancesOnClusterQuery = `
		SELECT * FROM instances
		WHERE cluster_id = $1
		ORDER BY id
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
	// DeleteClusterQuery removes an cluster by its name
	DeleteClusterQuery = `DELETE FROM clusters WHERE id=$1`

	// UpdateClusterQuery updates an existing cluster's mutable fields.
	UpdateClusterQuery = `
		UPDATE clusters SET
			owner = :owner,
			console_link = :console_link
		WHERE id = :id
	`

	// SelectProvidersOverviewQuery returns data about cloud providers with their account and cluster counts,
	// excluding those marked as "UNKNOWN" and not counting Terminated clusters
	SelectProvidersOverviewQuery = `
		SELECT
			a.provider,
			COUNT(DISTINCT a.name) AS account_count,
			COUNT(DISTINCT CASE WHEN c.status != 'Terminated' THEN c.id END) AS cluster_count
		FROM
			accounts a
		LEFT JOIN
			clusters c ON c.account_name = a.name
		WHERE
			a.provider != 'UNKNOWN'
		GROUP BY
			a.provider
		ORDER BY
			a.provider;
	` // TODO: review if this is needed here
	// CheckStatusQuery checks if the requested status exists on the DB
	CheckStatusQuery = `SELECT EXISTS (SELECT 1 FROM status WHERE value=$1)` // TODO: review if this is needed here
)
