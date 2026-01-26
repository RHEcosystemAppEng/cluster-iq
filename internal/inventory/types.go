package inventory

const (
	// ClusterTagKey string to identify to which cluster is the instance associated
	ClusterTagKey string = "kubernetes.io/cluster/"

	// Cluster actions
	ClusterPowerOnAction  = "PowerOn"
	ClusterPowerOffAction = "PowerOff"

	// Resource types
	ClusterResourceType  = "Cluster"
	InstanceResourceType = "Instance"
)
