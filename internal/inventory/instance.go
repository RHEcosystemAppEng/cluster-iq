package inventory

import (
	"fmt"
	"time"
)

// Instance model a cloud provider instance
// TODO: doc variables
type Instance struct {
	// Uniq Identifier of the instance
	ID string `db:"id" json:"id"`

	// Instance Name. In some Cloud Providers, the name is managed as a Tag
	Name string `db:"name" json:"name"`

	// Instance provider (public/private cloud provider)
	Provider CloudProvider `db:"provider" json:"provider"`

	// Instance type/size/flavour
	InstanceType string `db:"instance_type" json:"instanceType"`

	// Availability Zone in which the instance is running on
	AvailabilityZone string `db:"availability_zone" json:"availabilityZone"`

	// Instance Status
	State InstanceState `db:"state" json:"state"`

	// ClusterID
	ClusterID string `db:"cluster_id" json:"clusterID"`

	// Instance Tags as key-value array
	Tags []Tag `json:"tags"`

	// Instance Total Cost since created
	TotalCost float64 `db:"total_cost" json:"totalCost"`

	LastCostQueryTime time.Time `db:"last_cost_query_time" json:"lastCostQueryTime"`
}

// NewInstance returns a new Instance object
func NewInstance(id string, name string, provider CloudProvider, instanceType string, availabilityZone string, state InstanceState, clusterID string, tags []Tag, totalCost float64, lastCostQueryTime time.Time) *Instance {
	return &Instance{
		ID:                id,
		Name:              name,
		Provider:          provider,
		InstanceType:      instanceType,
		AvailabilityZone:  availabilityZone,
		State:             state,
		ClusterID:         clusterID,
		Tags:              tags,
		TotalCost:         totalCost,
		LastCostQueryTime: lastCostQueryTime,
	}
}

// AddTag adds a tag to an instance
func (i *Instance) AddTag(tag Tag) {
	i.Tags = append(i.Tags, tag)
}

// String as ToString func
func (i Instance) String() string {
	return fmt.Sprintf("%s(%s): [%s][%s][%s][%s][%s][TotalCost: %.2f][LastCostQueryTime: %s]",
		i.Name,
		i.ID,
		i.Provider,
		i.InstanceType,
		i.AvailabilityZone,
		i.State,
		i.ClusterID,
		i.TotalCost,
		i.LastCostQueryTime.Format(time.RFC3339),
	)
}

// PrintInstance prints Instance details
func (i Instance) PrintInstance() {
	fmt.Printf("\t\tInstance: %s\n", i.String())
}
