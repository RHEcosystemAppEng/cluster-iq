package inventory

import (
	"errors"
	"testing"
	"time"
)

var (
	id                string         = "01234"
	name              string         = "testInstance"
	provider          CloudProvider  = AWSProvider
	instanceType      string         = "medium"
	availabilityZone  string         = "eu-west-1a"
	status            InstanceStatus = Unknown
	clusterID         string         = "cluster-A01"
	lastScanTimestamp time.Time
	creationTimestamp time.Time
	age               int     = 1
	tags              []Tag   = []Tag{}
	totalCost         float64 = 0.0
)

func TestNewInstance(t *testing.T) {
	creationTimestamp, _ = time.Parse("2006-01-02", "2024-03-01")

	instance := NewInstance(id, name, provider, instanceType, availabilityZone, status, clusterID, tags, creationTimestamp, totalCost)
	if instance == nil {
		t.Errorf("Instance is null")
	}
}

func TestUpdateCosts(t *testing.T) {
	creationTimestamp, _ = time.Parse("2006-01-02", "2024-03-01")
	lastScanTimestamp, _ = time.Parse("2006-01-02", "2024-03-02")

	var tests = []struct {
		age               int
		totalCostExpected float64
		dailyCostExpected float64
		err               error
		expenses          []Expense
	}{
		{ // Case 1: Normal values
			10,
			34.0,
			3.4,
			nil,
			[]Expense{
				{
					"01234",
					24.0,
					time.Now(),
				},
				{
					"01234",
					10.0,
					time.Now(),
				},
			},
		},
		{ // Case 2: Negative amount
			2,
			6.0,
			3.0,
			nil,
			[]Expense{
				{
					"01234",
					-4.0,
					time.Now(),
				},
				{
					"01234",
					10.0,
					time.Now(),
				},
			},
		},
		{ // Case 3: Negative result
			2,
			0.0,
			0.0,
			ERR_INSTANCE_TOTAL_COST_LESS_ZERO,
			[]Expense{
				{
					"01234",
					-4.0,
					time.Now(),
				},
				{
					"01234",
					-12.0,
					time.Now(),
				},
			},
		},
		{ // Case 4: Age is zero
			0,
			6.0,
			0.0,
			ERR_INSTANCE_AGE_LESS_ZERO,
			[]Expense{
				{
					"01234",
					2.0,
					time.Now(),
				},
				{
					"01234",
					4.0,
					time.Now(),
				},
			},
		},
	}

	for _, test := range tests {
		instance := Instance{
			ID:                id,
			Name:              name,
			Provider:          provider,
			InstanceType:      instanceType,
			AvailabilityZone:  availabilityZone,
			Status:            status,
			ClusterID:         clusterID,
			LastScanTimestamp: lastScanTimestamp,
			CreationTimestamp: creationTimestamp,
			Age:               test.age,
			DailyCost:         0.0,
			TotalCost:         0.0,
			Tags:              tags,
			Expenses:          test.expenses,
		}

		result := instance.UpdateCosts()

		if !errors.Is(result, test.err) {
			t.Errorf("Expected errors mismatch. Have: %s ; Expected: %s", result, test.err)
		}
		if instance.TotalCost != test.totalCostExpected {
			t.Errorf("TotalCost is not correct. Have: %f ; Expected: %f", instance.TotalCost, test.totalCostExpected)
		}
		if instance.DailyCost != test.dailyCostExpected {
			t.Errorf("DailyCost is not correct. Have: %f ; Expected: %f", instance.DailyCost, test.dailyCostExpected)
		}
	}
}

func TestPrintInstance(t *testing.T) {
	var tests = []struct {
		instance *Instance
	}{
		{
			&Instance{
				ID:                id,
				Name:              name,
				Provider:          provider,
				InstanceType:      instanceType,
				AvailabilityZone:  availabilityZone,
				Status:            status,
				ClusterID:         clusterID,
				LastScanTimestamp: lastScanTimestamp,
				CreationTimestamp: creationTimestamp,
				Age:               2,
				DailyCost:         0.0,
				TotalCost:         totalCost,
				Tags:              tags,
			},
		},
		{nil},
	}

	for _, test := range tests {
		test.instance.PrintInstance()
	}
}

func TestAddTag(t *testing.T) {
	var tests = []struct {
		tags           []Tag
		expectedTagLen int
	}{
		{[]Tag{{Key: "KEY-1", Value: "-1"}, {Key: "KEY-2", Value: "-2"}}, 2}, // Case 1: Normal values
	}
	for _, test := range tests {
		instance := Instance{
			ID:                id,
			Name:              name,
			Provider:          provider,
			InstanceType:      instanceType,
			AvailabilityZone:  availabilityZone,
			Status:            status,
			ClusterID:         clusterID,
			LastScanTimestamp: lastScanTimestamp,
			CreationTimestamp: creationTimestamp,
			Age:               age,
			DailyCost:         0.0,
			TotalCost:         totalCost,
			Tags:              tags,
		}

		for _, tag := range test.tags {
			instance.AddTag(tag)
		}
		if len(instance.Tags) != test.expectedTagLen {
			t.Errorf("Tags lenght is not correct. Have: %d ; Expected: %d", len(instance.Tags), test.expectedTagLen)
		}
	}
}
