package cloudagent

// CloudExecutor interface defines the foundations for Executors. Executors are
// the implementation for connecting and sending orders to a specific cloud
// provider
type CloudExecutor interface {
	// Connect logs in into the cloud provider
	Connect() error
	// GetAccountName returns accounts name
	GetAccountName() string
	// SetRegion configure the cloud provider client for using a specific region
	SetRegion(string) error
	// PowerOffCluster triggers the power off process of a specific Cluster
	PowerOffCluster(instances []string)
	// PowerOnCluster triggers the power on process of a specific Cluster
	PowerOnCluster(instances []string)
}
