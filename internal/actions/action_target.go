package actions

// ActionTarget represents the target of an action in a cloud environment.
// It includes information about the account, region, cluster, and instances involved in the action.
type ActionTarget struct {
	// AccountID is the name of the cloud account associated with the target.
	AccountID string `db:"account_id"`

	// Region specifies the cloud region where the target resources are located.
	Region string `db:"region"`

	// ClusterID is the unique identifier of the cluster targeted by the action.
	ClusterID string `db:"cluster_id"`

	// Instances is a list of instance IDs associated with the target cluster.
	Instances []string `db:"instances"`
}

// NewActionTarget creates and returns a new instance of ActionTarget.
//
// Parameters:
// - AccountID: The name of the cloud account.
// - region: The cloud region where the resources are located.
// - clusterID: The unique identifier of the cluster.
// - instances: A slice of instance IDs associated with the cluster.
//
// Returns:
// - A pointer to a newly created ActionTarget instance.
func NewActionTarget(accountID string, region string, clusterID string, instances []string) *ActionTarget {
	return &ActionTarget{
		AccountID: accountID,
		Region:    region,
		ClusterID: clusterID,
		Instances: instances,
	}
}

// GetAccountID returns the name of the cloud account associated with the ActionTarget.
//
// Returns:
// - A string representing the account name.
func (at *ActionTarget) GetAccountID() string {
	return at.AccountID
}

// GetRegion returns the cloud region of the ActionTarget.
//
// Returns:
// - A string representing the region.
func (at *ActionTarget) GetRegion() string {
	return at.Region
}

// GetClusterID returns the unique identifier of the cluster targeted by the action.
//
// Returns:
// - A string representing the cluster ID.
func (at *ActionTarget) GetClusterID() string {
	return at.ClusterID
}

// GetInstances returns the list of instance IDs associated with the ActionTarget.
//
// Returns:
// - A slice of strings containing the instance IDs.
func (at *ActionTarget) GetInstances() []string {
	return at.Instances
}
