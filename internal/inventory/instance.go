package inventory

import (
	"errors"
	"fmt"
	"time"
)

// Errors for Instances
var (
	// TODO Remove!
	ERR_INSTANCE_TOTAL_COST_LESS_ZERO = errors.New("TotalCost of an instance cannot be less than zero")
	ERR_INSTANCE_DAILY_COST_LESS_ZERO = errors.New("DailyCost of an instance cannot be less than zero")
	ERR_INSTANCE_AGE_LESS_ZERO        = errors.New("Cannot recalculate costs if instance's Age is 0")
)

// TODO: Do we need the logic for calculating costs??

// Instance model a cloud provider instance
type Instance struct {
	InstanceID       string         `db:"instance_id"`       // Instance ID. Unique identifier of the instance in the cloud
	InstanceName     string         `db:"instance_name"`     // Instance Name. In some Cloud Providers, the name is managed as a Tag
	InstanceType     string         `db:"instance_type"`     // Instance type/size/flavour
	Provider         CloudProvider  `db:"provider"`          // Instance provider (public/private cloud provider)
	AvailabilityZone string         `db:"availability_zone"` // Availability Zone in which the instance is running on
	Status           ResourceStatus `db:"status"`            // Instance Status
	ClusterID        string         `db:"cluster_id"`        // ClusterID
	LastScanTS       time.Time      `db:"last_scan_ts"`      // Last scan timestamp of the instance
	CreatedAt        time.Time      `db:"created_at"`        // Timestamp when the instance was created
	Age              int            `db:"age"`               // Ammount of days since the instance was created

	// In-memory fields (not saved on DB)
	Cluster  *Cluster
	Tags     []Tag     `json:"tags"`     // Instance Tags as key-value array
	Expenses []Expense `json:"expenses"` // Expenses list associated to the instance
}

// NewInstance returns a new Instance object
func NewInstance(instanceID string, instanceName string, provider CloudProvider, instanceType string, availabilityZone string, status ResourceStatus, tags []Tag, creationTimestamp time.Time) *Instance {
	now := time.Now()
	age := calculateAge(creationTimestamp, now)

	return &Instance{
		InstanceID:       instanceID,
		InstanceName:     instanceName,
		Provider:         provider,
		InstanceType:     instanceType,
		AvailabilityZone: availabilityZone,
		Status:           status,
		ClusterID:        "",
		LastScanTS:       time.Time{},
		CreatedAt:        creationTimestamp,
		Age:              age,
		Cluster:          nil,
		Tags:             tags,
	}
}

// AddTag adds a tag to an instance
func (i *Instance) AddTag(tag Tag) {
	i.Tags = append(i.Tags, tag)
}

func (i *Instance) AddExpense(expense Expense) {
	i.Expenses = append(i.Expenses, expense)
}

// String as ToString func
func (i Instance) String() string {
	return fmt.Sprintf("(%s): [%s][%s][%s][%s][%s][%d]",
		i.InstanceName,
		i.Provider,
		i.InstanceType,
		i.AvailabilityZone,
		i.Status,
		i.ClusterID,
		len(i.Expenses),
	)
}

// PrintInstance prints Instance details
func (i Instance) PrintInstance() {
	fmt.Printf("\t\t\tInstance: %s\n", i.String())
}
