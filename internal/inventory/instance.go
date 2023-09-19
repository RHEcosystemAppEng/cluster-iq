package inventory

import (
	"fmt"
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

	// Region/Availability Zone in which the instance is running on
	Region string `db:"region" json:"region"`

	// Instance Status
	State InstanceState `db:"state" json:"state"`

	// ClusterName
	ClusterName string `db:"cluster_name" json:"clusterName"`

	// Instance Tags as key-value array
	Tags []Tag `json:"tags"`
}

// NewInstance returns a new Instance object
func NewInstance(id string, name string, provider CloudProvider, instanceType string, region string, state InstanceState, clusterName string, tags []Tag) *Instance {
	return &Instance{
		ID:           id,
		Name:         name,
		Provider:     provider,
		InstanceType: instanceType,
		Region:       region,
		State:        state,
		ClusterName:  clusterName,
		Tags:         tags,
	}
}

func (i Instance) String() string {
	return fmt.Sprintf("%s(%s): [%s][%s][%s][%s][%s]",
		i.Name,
		i.ID,
		i.Provider,
		i.InstanceType,
		i.Region,
		i.State,
		i.ClusterName,
	)
}

// PrintInstance prints Instance details
func (i Instance) PrintInstance() {
	fmt.Printf("\t\tInstance: %s\n", i.String())
}
