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
func ToInventoryInstanceList(dtos []InstanceDTORequest) *[]inventory.Instance {
	instances := make([]inventory.Instance, len(dtos))
	for i, dto := range dtos {
		instances[i] = *dto.ToInventoryInstance()
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
