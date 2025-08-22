package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TODO: comments
type InstanceDTORequest struct {
	InstanceID       string                   `json:"instanceID"`       // Instance's ID
	InstanceName     string                   `json:"instanceName"`     // Instance's Name
	InstanceType     string                   `json:"instanceType"`     // Instance's Type
	Provider         inventory.CloudProvider  `json:"provider"`         // Infrastructure provider identifier.
	AvailabilityZone string                   `json:"availabilityZone"` // The region of the infrastructure provider in which the cluster is deployed
	Status           inventory.ResourceStatus `json:"status"`           // Defines the status of the cluster if its infrastructure is running or not or it was removed
	ClusterID        string                   `json:"clusterID"`        // Account ID which this cluster belongs to
	LastScanTS       time.Time                `json:"lastScanTS"`       // Last scan timestamp of the cluster
	CreatedAt        time.Time                `json:"createdAt"`        // Timestamp when the cluster was created
	Age              int                      `json:"age"`              // Amount of days since the cluster was created
	Owner            string                   `json:"owner"`            // Cluster's owner
	Tags             TagDTORequestList        `json:"tags"`
}

func (i InstanceDTORequest) ToInventoryInstance() *inventory.Instance {
	instance := inventory.NewInstance(
		i.InstanceID,
		i.InstanceName,
		i.Provider,
		i.InstanceType,
		i.AvailabilityZone,
		i.Status,
		[]inventory.Tag{},
		i.CreatedAt,
	)

	instance.LastScanTS = i.LastScanTS
	instance.ClusterID = 0

	return instance
}

// TODO: comments
// InstanceDTORequestList represents the API Request containing a list of accounts.
type InstanceDTORequestList struct {
	Instances []InstanceDTORequest `json:"instances"` // List of accounts.
}

func (i InstanceDTORequestList) ToInventoryInstanceList() *[]inventory.Instance {
	var instances []inventory.Instance

	for _, instance := range i.Instances {
		instances = append(instances, *instance.ToInventoryInstance())
	}

	return &instances
}
