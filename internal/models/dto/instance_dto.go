package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// TODO: comments
type InstanceDTORequest struct {
	InstanceID        string                   `json:"instanceID"`
	InstanceName      string                   `json:"instanceName"`
	InstanceType      string                   `json:"instanceType"`
	Provider          inventory.Provider       `json:"provider"`
	AvailabilityZone  string                   `json:"availabilityZone"`
	Status            inventory.ResourceStatus `json:"status"`
	ClusterID         string                   `json:"clusterID"`
	LastScanTimestamp time.Time                `json:"lastScanTimestamp"`
	CreatedAt         time.Time                `json:"createdAt"`
	Age               int                      `json:"age"`
	Owner             string                   `json:"owner"`
	Tags              []TagDTORequest          `json:"tags"`
} // @name InstanceRequest

func (i InstanceDTORequest) ToInventoryInstance() *inventory.Instance {
	instance := inventory.NewInstance(
		i.InstanceID,
		i.InstanceName,
		i.Provider,
		i.InstanceType,
		i.AvailabilityZone,
		i.Status,
		*ToInventoryTagList(i.Tags),
		i.CreatedAt,
	)

	instance.LastScanTS = i.LastScanTimestamp
	instance.ClusterID = i.ClusterID

	return instance
}

func ToInventoryInstanceList(dtos []InstanceDTORequest) *[]inventory.Instance {
	instances := make([]inventory.Instance, len(dtos))
	for i, dto := range dtos {
		instances[i] = *dto.ToInventoryInstance()
	}

	return &instances
}

func ToInstanceDTORequest(instance inventory.Instance) *InstanceDTORequest {
	return &InstanceDTORequest{
		InstanceID:        instance.InstanceID,
		InstanceName:      instance.InstanceName,
		InstanceType:      instance.InstanceType,
		Provider:          instance.Provider,
		AvailabilityZone:  instance.AvailabilityZone,
		Status:            instance.Status,
		ClusterID:         instance.ClusterID,
		LastScanTimestamp: instance.LastScanTS,
		CreatedAt:         instance.CreatedAt,
		Age:               instance.Age,
		Owner:             inventory.GetOwnerFromTags(instance.Tags),
		Tags:              *ToTagDTORequestList(instance.Tags),
	}
}

func ToInstanceDTORequestList(instances []inventory.Instance) *[]InstanceDTORequest {
	instanceList := make([]InstanceDTORequest, len(instances))
	for i, instance := range instances {
		instanceList[i] = *ToInstanceDTORequest(instance)
	}

	return &instanceList
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
	LastScanTimestamp     time.Time                `json:"lastScanTimestamp"`
	CreatedAt             time.Time                `json:"creationTimestamp"`
	Age                   int                      `json:"age"`
	Owner                 string                   `json:"owner"`
	TotalCost             float64                  `json:"totalCost"`
	Last15DaysCost        float64                  `json:"last15DaysCost"`
	LastMonthCost         float64                  `json:"lastMonthCost"`
	CurrentMonthSoFarCost float64                  `json:"currentMonthSoFarCost"`
	Tags                  []TagDTOResponse         `json:"tags"`
} // @name InstanceResponse

func (i *InstanceDTOResponse) ToInventoryInstance() *inventory.Instance {
	return &inventory.Instance{
		InstanceID:       i.InstanceID,
		InstanceName:     i.InstanceName,
		InstanceType:     i.InstanceType,
		Provider:         i.Provider,
		AvailabilityZone: i.AvailabilityZone,
		Status:           i.Status,
		ClusterID:        i.ClusterID,
		LastScanTS:       i.LastScanTimestamp,
		CreatedAt:        i.CreatedAt,
		Age:              i.Age,
	}
}
