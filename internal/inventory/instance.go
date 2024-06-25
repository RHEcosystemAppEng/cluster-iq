package inventory

import (
	"errors"
	"fmt"
	"time"
)

// Errors for Instances
var (
	ERR_INSTANCE_TOTAL_COST_LESS_ZERO = errors.New("TotalCost of an instance cannot be less than zero")
	ERR_INSTANCE_DAILY_COST_LESS_ZERO = errors.New("DailyCost of an instance cannot be less than zero")
	ERR_INSTANCE_AGE_LESS_ZERO        = errors.New("Cannot recalculate costs if instance's Age is 0")
)

// Instance model a cloud provider instance
type Instance struct {
	// Uniq Identifier of the instance
	ID string `db:"id" json:"id"`

	// Instance Name. In some Cloud Providers, the name is managed as a Tag
	Name string `db:"name" json:"name"`

	// Instance provider (public/private cloud provider)
	Provider CloudProvider `db:"provider" json:"provider"`

	// Instance type/size/flavour
	InstanceType string `db:"instance_type" json:"instanceType"`

	// Availability Zone in which the instance is running on
	AvailabilityZone string `db:"availability_zone" json:"availabilityZone"`

	// Instance Status
	Status InstanceStatus `db:"status" json:"status"`

	// ClusterID
	ClusterID string `db:"cluster_id" json:"clusterID"`

	// Last scan timestamp of the instance
	LastScanTimestamp time.Time `db:"last_scan_timestamp" json:"lastScanTimestamp"`

	// Timestamp when the instance was created
	CreationTimestamp time.Time `db:"creation_timestamp" json:"creationTimestamp"`

	// Ammount of days since the instance was created
	Age int `db:"age" json:"age"`

	// Daily cost (US Dollars) estimated based on total cost and age of the instance
	DailyCost float64 `db:"daily_cost" json:"dailyCost"`

	// Total cost (US Dollars) accumulated since ClusterIQ is scanning
	TotalCost float64 `db:"total_cost" json:"totalCost"`

	// Instance Tags as key-value array
	Tags []Tag `json:"tags"`

	// Expenses list associated to the instance
	Expenses []Expense `json:"expenses"`
}

// NewInstance returns a new Instance object
func NewInstance(id string, name string, provider CloudProvider, instanceType string, availabilityZone string, status InstanceStatus, clusterID string, tags []Tag, creationTimestamp time.Time, totalCost float64) *Instance {
	now := time.Now()
	age := calculateAge(creationTimestamp, now)

	return &Instance{
		ID:                id,
		Name:              name,
		Provider:          provider,
		InstanceType:      instanceType,
		AvailabilityZone:  availabilityZone,
		Status:            status,
		ClusterID:         clusterID,
		LastScanTimestamp: now,
		CreationTimestamp: creationTimestamp,
		Age:               age,
		DailyCost:         0.0,
		TotalCost:         totalCost,
		Tags:              tags,
	}
}

// SetTotalCost sets the TotalCost of an instance and recalculates the rest of costs
func (i *Instance) calculateTotalCost() error {
	var totalCost float64 = 0.0
	for _, expense := range i.Expenses {
		totalCost += expense.Amount
	}

	if totalCost < 0 {
		return ERR_INSTANCE_TOTAL_COST_LESS_ZERO
	}

	i.TotalCost = totalCost
	return nil
}

// calculateDailyCost calculates and retruns the average Amount of the instance per day
func (i *Instance) calculateDailyCost() error {
	var dailyCost float64 = 0.0

	if i.Age <= 0 {
		return ERR_INSTANCE_AGE_LESS_ZERO
	}
	dailyCost = i.TotalCost / float64(i.Age)

	if dailyCost < 0 {
		return ERR_INSTANCE_DAILY_COST_LESS_ZERO
	}

	i.DailyCost = dailyCost
	return nil
}

// UpdateCosts updates the totalCost of the instance using the instance age and the DailyCost
func (i *Instance) UpdateCosts() error {
	if err := i.calculateTotalCost(); err != nil {
		return err
	}

	if err := i.calculateDailyCost(); err != nil {
		return err
	}

	return nil
}

// AddTag adds a tag to an instance
func (i *Instance) AddTag(tag Tag) {
	i.Tags = append(i.Tags, tag)
}

// PrintInstance prints Instance details
func (i *Instance) PrintInstance() {
	if i == nil {
		return
	}
	if str, err := JSONMarshal(i); err != nil {
		fmt.Printf("\t\tInstance: %s\n", str)
	} else {
		fmt.Println(i)
	}
}
