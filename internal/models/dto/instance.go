package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TODO: comments
type InstanceDTOResponse struct {
	InstanceID            string                   `json:"instanceID"`        // Instance's ID
	InstanceName          string                   `json:"instanceName"`      // Instance's Name
	InstanceType          string                   `json:"instanceType"`      // Instance's Type
	Provider              inventory.CloudProvider  `json:"provider"`          // Infrastructure provider identifier.
	AvailabilityZone      string                   `json:"availabilityZone"`  // The region of the infrastructure provider in which the cluster is deployed
	Status                inventory.ResourceStatus `json:"status"`            // Defines the status of the cluster if its infrastructure is running or not or it was removed
	ClusterID             string                   `json:"clusterID"`         // Account ID which this cluster belongs to
	ClusterName           string                   `json:"clusterName"`       // Account ID which this cluster belongs to
	LastScanTS            time.Time                `json:"lastScanTimestamp"` // Last scan timestamp of the cluster
	CreatedAt             time.Time                `json:"creationTimestamp"` // Timestamp when the cluster was created
	Age                   int                      `json:"age"`               // Amount of days since the cluster was created
	Owner                 string                   `json:"owner"`             // Cluster's owner
	TotalCost             float64                  `json:"totalCost"`
	Last15DaysCost        float64                  `json:"last15DaysCost"`
	LastMonthCost         float64                  `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64                  `json:"currentMonthSoFarCost"`
}

// TODO: comments
// InstanceDTOResponseList represents the API response containing a list of accounts.
type InstanceDTOResponseList struct {
	Count     int                   `json:"count,omitempty"` // Number of accounts, omitted if empty.
	Instances []InstanceDTOResponse `json:"instances"`       // List of accounts.
}

// TODO: comments
// NewInstanceDTOResponseList creates a new InstanceDTOResponseList instance.
// It ensures that an empty array is returned if the input account list is empty.
//
// Parameters:
// - accounts: A slice of inventory.Account.
//
// Returns:
// - A pointer to an InstanceDTOResponseList.
func NewInstanceDTOResponseList(instances []InstanceDTOResponse) *InstanceDTOResponseList {
	response := InstanceDTOResponseList{Instances: instances}

	// Count only set list length > 0
	if count := len(instances); count > 0 {
		response.Count = count
	}

	return &response
}

// TODO: comments
type InstanceDTORequest struct {
	InstanceID       string                   `json:"instanceID"`        // Instance's ID
	InstanceName     string                   `json:"instanceName"`      // Instance's Name
	InstanceType     string                   `json:"instanceType"`      // Instance's Type
	Provider         inventory.CloudProvider  `json:"provider"`          // Infrastructure provider identifier.
	AvailabilityZone string                   `json:"availabilityZone"`  // The region of the infrastructure provider in which the cluster is deployed
	Status           inventory.ResourceStatus `json:"status"`            // Defines the status of the cluster if its infrastructure is running or not or it was removed
	ClusterID        string                   `json:"clusterID"`         // Account ID which this cluster belongs to
	LastScanTS       time.Time                `json:"lastScanTimestamp"` // Last scan timestamp of the cluster
	CreatedAt        time.Time                `json:"creationTimestamp"` // Timestamp when the cluster was created
	Age              int                      `json:"age"`               // Amount of days since the cluster was created
	Owner            string                   `json:"owner"`             // Cluster's owner
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
	instance.ClusterID = i.ClusterID

	return instance
}

// TODO: comments
// InstanceDTORequestList represents the API Request containing a list of accounts.
type InstanceDTORequestList struct {
	Instances []InstanceDTORequest `json:"instances"` // List of accounts.
}
