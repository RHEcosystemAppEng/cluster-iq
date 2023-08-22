package inventory

import "fmt"

// Cluster stores the cluster properties in the stock
// TODO: doc variables
type Cluster struct {
	Name        string        `redis:"name" json:"name"`
	AccountName string        `redis:"accountName" json:"accountName"`
	Provider    CloudProvider `redis:"provider" json:"provider"`
	Status      InstanceState `redis:"status" json:"status"`
	Region      string        `redis:"region" json:"region"`
	ConsoleLink string        `redis:"consoleLink" json:"consoleLink"`
	Instances   []Instance    `redis:"instances" json:"instances"`
}

// NewCluster creates a new cluster instance
func NewCluster(name string, accountName string, provider CloudProvider, region string, consoleLink string) Cluster {
	return Cluster{
		Name:        name,
		AccountName: accountName,
		Provider:    provider,
		Status:      Unknown,
		Region:      region,
		ConsoleLink: consoleLink,
		Instances:   make([]Instance, 0),
	}
}

// TODO: doc
func (c Cluster) checkStatus(state InstanceState) bool {
	instancesNum := len(c.Instances)
	count := 0
	for _, instance := range c.Instances {
		if instance.State == state {
			count++
		}
	}

	if count >= (instancesNum / 2) {
		return true
	}

	return false
}

// TODO: doc
func (c Cluster) isClusterStopped() bool {
	return c.checkStatus(Stopped)
}

// TODO: doc
func (c Cluster) isClusterRunning() bool {
	return c.checkStatus(Running)
}

// UpdateStatus checks the instance list and updates the cluster status based
// on its instances status. If half or more of the nodes are on the same
// status, cluster will have the same
func (c *Cluster) UpdateStatus() {
	if c.isClusterRunning() {
		c.Status = Running
	} else if c.isClusterStopped() {
		c.Status = Stopped
	} else {
		c.Status = Unknown
	}
}

// AddInstance add a new instance to a cluster
func (c *Cluster) AddInstance(instance Instance) {
	c.Instances = append(c.Instances, instance)
	c.UpdateStatus()
}

// PrintCluster prints cluster info
func (c Cluster) PrintCluster() {
	fmt.Printf("\tCluster: %s -- [%s][%s](Instances: %d)\n", c.Name, c.ConsoleLink, c.AccountName, len(c.Instances))
	for _, instance := range c.Instances {
		instance.PrintInstance()
	}
	fmt.Printf("\n")
}
