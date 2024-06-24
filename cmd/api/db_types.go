package main

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

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

	// Availability Zone in which the instance is running on
	AvailabilityZone string `db:"availability_zone"`

	// Instance Status
	State inventory.InstanceState `db:"state"`

	// ClusterID
	ClusterID string `db:"cluster_id"`

	// instance Tags
	TagKey   string `db:"key"`
	TagValue string `db:"value"`

	// InstanceID from Tags table (not needed, defined just for parsing the join result)
	InstanceID string `db:"instance_id"`

	// Instance TotalCost
	TotalCost float64 `db:"total_cost"`

	// Last cost query time
	LastCostQueryTime time.Time `db:"last_cost_query_time"`
}
