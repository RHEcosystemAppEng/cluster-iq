package dto

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"time"
)

// Instance model a cloud provider instance
type Instance struct {
	// Uniq Identifier of the instance
	ID string `json:"id"`

	// Instance Name. In some Cloud Providers, the name is managed as a Tag
	Name string `json:"name"`

	// Instance provider (public/private cloud provider)
	//TODO
	Provider inventory.CloudProvider `json:"provider"`

	// Instance type/size/flavour
	InstanceType string `json:"instanceType"`

	// Availability Zone in which the instance is running on
	AvailabilityZone string `json:"availabilityZone"`

	// Instance Status
	//TODO
	Status inventory.InstanceStatus `json:"status"`

	// ClusterID
	ClusterID string `json:"clusterID"`

	// Last scan timestamp of the instance
	LastScanTimestamp time.Time `json:"lastScanTimestamp"`

	// Timestamp when the instance was created
	CreationTimestamp time.Time `json:"creationTimestamp"`

	// Amount of days since the instance was created
	Age int `json:"age"`

	// Daily cost (US Dollars) estimated based on total cost and age of the instance
	DailyCost float64 `json:"dailyCost"`

	// Total cost (US Dollars) accumulated since ClusterIQ is scanning
	TotalCost float64 `json:"totalCost"`

	// Instance Tags as key-value array
	Tags []Tag `json:"tags"`

	// Expenses list associated to the instance
	// TODO, remove?????
	Expenses []inventory.Expense `json:"expenses"`
}
