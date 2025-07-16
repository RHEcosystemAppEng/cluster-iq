package inventory

const (
	// ClusterTagKey string to identify to which cluster is the instance associated
	ClusterTagKey string = "kubernetes.io/cluster/"
)

const (
	// Cluster actions
	ClusterPowerOnAction  = "PowerOn"
	ClusterPowerOffAction = "PowerOff"

	// Resource types
	ClusterResourceType  = "cluster"
	InstanceResourceType = "instance"
)
