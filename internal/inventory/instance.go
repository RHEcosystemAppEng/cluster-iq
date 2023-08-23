package inventory

import (
	"fmt"
)

// Instance model a cloud provider instance
// TODO: doc variables
type Instance struct {
	ID           string        `redis:"id" json:"id"`
	Name         string        `redis:"name" json:"name"`
	Region       string        `redis:"region" json:"region"`
	InstanceType string        `redis:"instanceType" json:"instanceType"`
	State        InstanceState `redis:"state" json:"state"`
	Provider     CloudProvider `redis:"provider" json:"provider"`
	Tags         []Tag         `redis:"tags" json:"tags"`
}

// NewInstance returns a new Instance object
func NewInstance(
	id string,
	name string,
	region string,
	instanceType string,
	state InstanceState,
	provider CloudProvider,
	tags []Tag,
) Instance {
	return Instance{
		ID:           id,
		Name:         name,
		Region:       region,
		InstanceType: instanceType,
		State:        state,
		Provider:     provider,
		Tags:         tags,
	}
}

func (i Instance) String() string {
	return fmt.Sprintf("%s(%s): [%s][%s][%s][%s]",
		i.Name,
		i.ID,
		i.Provider,
		i.Region,
		i.State,
		i.InstanceType,
	)
}

// PrintInstance prints Instance details
func (i Instance) PrintInstance() {
	fmt.Printf("\t\tInstance: %s\n", i.String())
}
