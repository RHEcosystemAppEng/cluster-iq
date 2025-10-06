package repositories

const (
	InstancesTable = "instances"
	// View for SELECT operations on Instances
	SelectInstancesFullView = "instances_full_view"
	// View for SELECT operations on Instances including tags
	SelectInstancesFullWithTagsView = "instances_full_view_with_tags"
	// View for SELECT operations for Instances pending on Expense Update
	SelectInstancesPendingExpenseUpdate = "instances_pending_expense_update"

	// Run 'check_terminated_instances()' SQL function
	UpdateTerminatedInstancesQuery = `SELECT check_terminated_instances()`
	// Run 'check_terminated_clusters()' SQL function
	UpdateTerminatedClustersQuery = `SELECT check_terminated_clusters()`
	// UpdateInstanceStatus updates the status of a  set of instances based on their clusterID
	UpdateStatusInstancesByClusterIDQuery = `UPDATE instances SET status=$1 WHERE cluster_id=$2`

	// InsertInstancesQuery inserts into a new instance in its table
	InsertInstancesQuery = `
		INSERT INTO instances (
			instance_id,
			instance_name,
			provider,
			instance_type,
			availability_zone,
			status,
			cluster_id,
			last_scan_ts,
			created_at,
			age
		) VALUES (
			:instance_id,
			:instance_name,
			:provider,
			:instance_type,
			:availability_zone,
			:status,
			(SELECT id FROM clusters WHERE cluster_id = :cluster_id),
			:last_scan_ts,
			:created_at,
			:age
		) ON CONFLICT (instance_id, cluster_id) DO UPDATE SET
			instance_name = EXCLUDED.instance_name,
			provider = EXCLUDED.provider,
			instance_type = EXCLUDED.instance_type,
			availability_zone = EXCLUDED.availability_zone,
			status = EXCLUDED.status,
			cluster_id = EXCLUDED.cluster_id,
			last_scan_ts = EXCLUDED.last_scan_ts,
			created_at = EXCLUDED.created_at,
			age = EXCLUDED.age
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
			(SELECT id FROM instances WHERE instance_id = :instance_id)
		) ON CONFLICT (key, instance_id) DO UPDATE SET
			value = EXCLUDED.value
	`
)
