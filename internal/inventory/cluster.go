package inventory

import "fmt"

const minInstances int = 3

// Cluster is the object to store Openshift Clusters and its properties
type Cluster struct {
	// Cluster's Name. Must be unique per Account
	Name string `redis:"name" json:"name"`

	// Infrastructure provider identifier.
	Provider CloudProvider `redis:"provider" json:"provider"`

	// Defines the status of the cluster if its infrastructure is running or not
	Status InstanceState `redis:"status" json:"status"`

	// The region of the infrastructure provider in which the cluster is deployed
	Region string `redis:"region" json:"region"`

	// Openshift Console URL. Might not be accesible if its protected
	ConsoleLink string `redis:"consoleLink" json:"consoleLink"`

	// Cluster's instance (nodes) lists
	Instances []Instance `redis:"instances" json:"instances"`
}

// NewCluster creates a new cluster instance
func NewCluster(name string, provider CloudProvider, region string, consoleLink string) Cluster {
	return Cluster{
		Name:        name,
		Provider:    provider,
		Status:      Unknown,
		Region:      region,
		ConsoleLink: consoleLink,
		Instances:   make([]Instance, 0),
	}
}

// isClusterStopped checks if the Cluster is Stopped
func (c Cluster) isClusterStopped() bool {
	if c.Status == Stopped {
		return true
	}
	return false
}

// isClusterRunning checks if the Cluster is Running
func (c Cluster) isClusterRunning() bool {
	if c.Status == Running {
		return true
	}
	return false
}

// UpdateStatus evaluate the status of the cluster checking how many of the
// nodes are in Running or Stopped status. As Openshift needs at lease 3 nodes
// running to be considered correctly Running (3 master nodes), but we cant'
// figure out which Instance is a master node, if at least 3 of the Cluster
// instances are running, Cluster will be considered as Running also.
// If the instances count is less than minInstances, Cluster would be
// considered as Unknown status
// TODO: find out a more trustable approach to evaluate cluster status
func (c *Cluster) UpdateStatus() {
	instancesNum := len(c.Instances)

	// Check minimun instances
	if instancesNum < minInstances {
		c.Status = Unknown
		return
	}

	count := 0
	for _, instance := range c.Instances {
		if instance.State == Running {
			count++
		}
		if count >= minInstances {
			c.Status = Running
			return
		}
	}

	c.Status = Stopped
}

// AddInstance add a new instance to a cluster
func (c *Cluster) AddInstance(instance Instance) {
	c.Instances = append(c.Instances, instance)
	c.UpdateStatus()
}

// PrintCluster prints cluster info
func (c Cluster) PrintCluster() {
	fmt.Printf("\tCluster: %s -- [%s](Instances: %d)\n", c.Name, c.ConsoleLink, len(c.Instances))
	for _, instance := range c.Instances {
		instance.PrintInstance()
	}
	fmt.Printf("\n")
}
