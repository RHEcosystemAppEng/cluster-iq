package repositories

const (
	// DeleteTagsQuery removes a Tag by its key and instance reference
	DeleteTagsQuery = `DELETE FROM tags WHERE instance_id=$1`
	// DeleteInstanceQuery removes an instance by its ID
	DeleteInstanceQuery = `DELETE FROM instances WHERE id=$1`
	// Run 'check_terminated_instances()' SQL function
	UpdateTerminatedInstancesQuery = `SELECT check_terminated_instances()`
	// Run 'check_terminated_clusters()' SQL function
	UpdateTerminatedClustersQuery = `SELECT check_terminated_clusters()`
	// UpdateInstanceStatus updates the status of a  set of instances based on their clusterID
	UpdateStatusInstancesByClusterIDQuery = `UPDATE instances SET status=$1 WHERE cluster_id=$2`
	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = `
		SELECT * FROM instances
	`
	// SelectInstancesOverview returns the total count of all instances
	SelectInstancesOverview = `
		SELECT
			COUNT(CASE WHEN status = 'Running' THEN 1 END) AS running,
			COUNT(CASE WHEN status = 'Stopped' THEN 1 END) AS stopped,
			COUNT(CASE WHEN status = 'Terminated' THEN 1 END) AS archived
		FROM instances;
	`

	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = `
		SELECT * FROM instances
		WHERE id = $1
	`
	// SelectTagsByInstanceIDQuery returns all tags for a single instance
	SelectTagsByInstanceIDQuery = `
		SELECT * FROM tags WHERE instance_id = $1
	`
	// SelectTagsByInstanceIDsQuery returns all tags for a list of instance IDs
	SelectTagsByInstanceIDsQuery = `
		SELECT * FROM tags WHERE instance_id IN (?)
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
)
