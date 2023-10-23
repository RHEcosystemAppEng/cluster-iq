package main

import "github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"

// InstanceDB is an intermediate struct to map the Instances and its tags into inventory.Instance objects
type InstanceDB struct {
	// Uniq Identifier of the instance
	ID string `db:"id"`

	// Instance Name. In some Cloud Providers, the name is managed as a Tag
	Name string `db:"name"`

	// Instance provider (public/private cloud provider)
	Provider inventory.CloudProvider `db:"provider"`

	// Instance type/size/flavour
	InstanceType string `db:"instance_type"`

	// Region/Availability Zone in which the instance is running on
	Region string `db:"region"`

	// Instance Status
	State inventory.InstanceState `db:"state"`

	// ClusterName
	ClusterName string `db:"cluster_name"`

	// instance Tags
	TagKey   string `db:"key"`
	TagValue string `db:"value"`

	// InstanceID from Tags table (not needed, defined just for parsing the join result)
	InstanceID string `db:"instance_id"`
}
