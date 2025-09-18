package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TODO: comments
type InstanceDTORequest struct {
	InstanceID       string                   `json:"instanceID"`
	InstanceName     string                   `json:"instanceName"`
	InstanceType     string                   `json:"instanceType"`
	Provider         inventory.Provider       `json:"provider"`
	AvailabilityZone string                   `json:"availabilityZone"`
	Status           inventory.ResourceStatus `json:"status"`
	ClusterID        string                   `json:"clusterID"`
	LastScanTS       time.Time                `json:"lastScanTS"`
	CreatedAt        time.Time                `json:"createdAt"`
	Age              int                      `json:"age"`
	Owner            string                   `json:"owner"`
	Tags             []TagDTORequest          `json:"tags"`
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
	instance.ClusterID = ""

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

// TODO: comments
type InstanceDTOResponse struct {
	InstanceID            string                   `json:"instanceID"`
	InstanceName          string                   `json:"instanceName"`
	InstanceType          string                   `json:"instanceType"`
	Provider              inventory.Provider       `json:"provider"`
	AvailabilityZone      string                   `json:"availabilityZone"`
	Status                inventory.ResourceStatus `json:"status"`
	ClusterID             string                   `json:"clusterID"`
	ClusterName           string                   `json:"clusterName"`
	LastScanTS            time.Time                `json:"lastScanTimestamp"`
	CreatedAt             time.Time                `json:"creationTimestamp"`
	Age                   int                      `json:"age"`
	Owner                 string                   `json:"owner"`
	TotalCost             float64                  `json:"totalCost"`
	Last15DaysCost        float64                  `json:"last15DaysCost"`
	LastMonthCost         float64                  `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64                  `json:"currentMonthSoFarCost"`
	Tags                  []TagDTOResponse         `json:"tags"`
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
