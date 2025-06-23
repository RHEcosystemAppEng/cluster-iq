package inventory

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// TestNewInstance verifies that NewInstance returns a correctly initialized instance
func TestNewInstance(t *testing.T) {
	now := time.Now().Add(-48 * time.Hour)
	inst := NewInstance("i-123", "test-instance", AWSProvider, "t3.micro", "us-east-1a", Running, "cluster-1", nil, now)

	if inst.ID != "i-123" {
		t.Errorf("expected ID i-123, got %s", inst.ID)
	}
	if inst.Age != 2 {
		t.Errorf("expected age 2, got %d", inst.Age)
	}
	if inst.TotalCost != 0.0 || inst.DailyCost != 0.0 {
		t.Error("expected initial costs to be 0")
	}
}

// TestCalculateTotalCost_Success verifies that total cost is correctly aggregated
func TestCalculateTotalCost_Success(t *testing.T) {
	i := Instance{
		Expenses: []Expense{
			{Amount: 2.5},
			{Amount: 3.5},
		},
	}
	err := i.calculateTotalCost()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i.TotalCost != 6.0 {
		t.Errorf("expected total cost 6.0, got %f", i.TotalCost)
	}
}

// TestCalculateTotalCost_ErrorNegative verifies that total cost fails when resulting value is negative
func TestCalculateTotalCost_ErrorNegative(t *testing.T) {
	i := Instance{
		Expenses: []Expense{{Amount: -10}},
	}
	err := i.calculateTotalCost()
	if !errors.Is(err, ERR_INSTANCE_TOTAL_COST_LESS_ZERO) {
		t.Errorf("expected ERR_INSTANCE_TOTAL_COST_LESS_ZERO, got %v", err)
	}
}

// TestCalculateDailyCost_Success verifies that daily cost is computed correctly
func TestCalculateDailyCost_Success(t *testing.T) {
	i := Instance{TotalCost: 10.0, Age: 5}
	err := i.calculateDailyCost()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i.DailyCost != 2.0 {
		t.Errorf("expected daily cost 2.0, got %f", i.DailyCost)
	}
}

// TestCalculateDailyCost_ZeroAge verifies error when age is zero
func TestCalculateDailyCost_ZeroAge(t *testing.T) {
	i := Instance{TotalCost: 10.0, Age: 0}
	err := i.calculateDailyCost()
	if !errors.Is(err, ERR_INSTANCE_AGE_LESS_ZERO) {
		t.Errorf("expected ERR_INSTANCE_AGE_LESS_ZERO, got %v", err)
	}
}

// TestCalculateDailyCost_Negative verifies error when daily cost would be negative
func TestCalculateDailyCost_Negative(t *testing.T) {
	i := Instance{TotalCost: -10.0, Age: 5}
	err := i.calculateDailyCost()
	if !errors.Is(err, ERR_INSTANCE_DAILY_COST_LESS_ZERO) {
		t.Errorf("expected ERR_INSTANCE_DAILY_COST_LESS_ZERO, got %v", err)
	}
}

// TestUpdateCosts_Success verifies UpdateCosts runs correctly when no errors are present
func TestUpdateCosts_Success(t *testing.T) {
	i := Instance{
		Age:      2,
		Expenses: []Expense{{Amount: 8.0}},
	}
	err := i.UpdateCosts()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i.TotalCost != 8.0 || i.DailyCost != 4.0 {
		t.Errorf("unexpected costs: total %f, daily %f", i.TotalCost, i.DailyCost)
	}
}

// TestUpdateCosts_TotalCostError verifies UpdateCosts fails on invalid total cost
func TestUpdateCosts_TotalCostError(t *testing.T) {
	i := Instance{
		Age:      2,
		Expenses: []Expense{{Amount: -1}},
	}
	err := i.UpdateCosts()
	if !errors.Is(err, ERR_INSTANCE_TOTAL_COST_LESS_ZERO) {
		t.Errorf("expected ERR_INSTANCE_TOTAL_COST_LESS_ZERO, got %v", err)
	}
}

// TestUpdateCosts_DailyCostError verifies UpdateCosts fails on invalid daily cost
func TestUpdateCosts_DailyCostError(t *testing.T) {
	i := Instance{
		Age:      0,
		Expenses: []Expense{{Amount: 10}},
	}
	err := i.UpdateCosts()
	if !errors.Is(err, ERR_INSTANCE_AGE_LESS_ZERO) {
		t.Errorf("expected ERR_INSTANCE_AGE_LESS_ZERO, got %v", err)
	}
}

// TestAddTag verifies that a new tag is appended to the tag list
func TestAddTag(t *testing.T) {
	i := Instance{}
	tag := Tag{Key: "env", Value: "prod"}
	i.AddTag(tag)
	if len(i.Tags) != 1 || i.Tags[0] != tag {
		t.Errorf("tag was not added correctly")
	}
}

// TestInstance_String verifies String method returns expected format
func TestInstance_String(t *testing.T) {
	i := Instance{
		ID:               "i-123",
		Name:             "test",
		Provider:         AWSProvider,
		InstanceType:     "t2.micro",
		AvailabilityZone: "us-east-1a",
		Status:           Running,
		ClusterID:        "cluster-x",
		Expenses:         []Expense{{Amount: 5}},
	}

	str := i.String()
	if !(strings.Contains(str, "test") && strings.Contains(str, "AWS") && strings.Contains(str, "t2.micro")) {
		t.Errorf("unexpected output from String(): %s", str)
	}
}

// TestPrintInstance verifies PrintInstance runs without panic
func TestPrintInstance(t *testing.T) {
	i := Instance{ID: "i-456", Name: "node1"}
	i.PrintInstance()
}
