package inventory

import (
	"errors"
	"fmt"
	"time"
)

// Errors for Instances
var (
	// Error when creating a new Instance without an InstanceID
	ErrMissingInstanceIDCreation = errors.New("cannot create an Instance without InstanceID")
	// Error when the total cost of an instance is less than 0
	ErrInstanceTotalCostLessZero = errors.New("total cost of an instance cannot be less than 0.0")
	// Error when the instance daily cost is less than 0
	ErrInstanceDailyCostLessZero = errors.New("daily cost of an instance cannot be less than 0.0")
	// Error when the calculated Age for an instance is less than 0
	ErrInstanceAgeLessZero = errors.New("instance age cannot be less than 0")
	// Error when adding a tag without Key to an Instance
	ErrAddingTagWithoutKey = errors.New("cannot add keyless tags")
	// Error when adding expenses without a valid amount to an instance
	ErrAddingExpenseWithWrongAmount = errors.New("cannot add an expense with negative amount")
)

// Instance model a cloud provider instance
type Instance struct {
	// InstanceID is the unique identifier of the instance defined by the Provider
	InstanceID string `db:"instance_id"`

	// InstanceName. In some Cloud Providers, the name is managed as a Tag
	InstanceName string `db:"instance_name"`

	// Instance type/size/flavour
	InstanceType string `db:"instance_type"`

	// Provider identifies the cloud/infrastructure provider.
	Provider Provider `db:"provider"`

	// AvailabilityZone in which the instance is running on
	AvailabilityZone string `db:"availability_zone"`

	// Status defines the status of the instance if it's running or not or it was removed (Terminated).
	Status ResourceStatus `db:"status"`

	// ClusterID which the instance is part of
	ClusterID string `db:"cluster_id"`

	// LastScanTimestamp is the timestamp when the instance was scanned for the last time.
	LastScanTimestamp time.Time `db:"last_scan_ts"`

	// CreatedAt is the timestamp when the instance was created (from the inventory point of view, not from the provider).
	CreatedAt time.Time `db:"created_at"`

	// Age is the amount of days since the cluster was created.
	Age int `db:"age"`

	// In-memory fields (no saved on DB)
	// ===========================================================================

	// Instance Tags as key-value array
	Tags []Tag

	// Expenses list associated to the instance
	Expenses []Expense
}

// NewInstance returns a new Instance object
func NewInstance(instanceID string, instanceName string, provider Provider, instanceType string, availabilityZone string, status ResourceStatus, tags []Tag, creationTimestamp time.Time) (*Instance, error) {
	if instanceID == "" {
		return nil, ErrMissingInstanceIDCreation
	}

	now := time.Now()
	age := calculateAge(creationTimestamp, now)

	return &Instance{
		InstanceID:        instanceID,
		InstanceName:      instanceName,
		Provider:          provider,
		InstanceType:      instanceType,
		AvailabilityZone:  availabilityZone,
		Status:            status,
		ClusterID:         "",
		LastScanTimestamp: now,
		CreatedAt:         creationTimestamp,
		Age:               age,
		Tags:              tags,
		Expenses:          make([]Expense, 0),
	}, nil
}

// AddTag adds a tag to an instance
func (i *Instance) AddTag(tag Tag) error {
	if tag.Key == "" {
		return ErrAddingTagWithoutKey
	}
	i.Tags = append(i.Tags, tag)

	return nil
}

func (i *Instance) AddExpense(expense *Expense) error {
	if expense.Amount < 0 {
		return ErrAddingExpenseWithWrongAmount
	}

	// Asigning the new instanceID and adding to the list
	expense.InstanceID = i.InstanceID
	i.Expenses = append(i.Expenses, *expense)

	return nil
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
