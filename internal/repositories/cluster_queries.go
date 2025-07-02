package repositories

const (
	// SelectClustersQuery returns every cluster in the inventory ordered by Name
	SelectClustersQuery = `
		SELECT * FROM clusters
	`
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
	// DeleteClusterQuery removes an cluster by its name
	DeleteClusterQuery = `DELETE FROM clusters WHERE id=$1`
)
