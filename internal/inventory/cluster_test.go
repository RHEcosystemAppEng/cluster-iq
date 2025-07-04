package inventory

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGenerateClusterID tests GenerateClusterID with valid and invalid inputs
func TestGenerateClusterID(t *testing.T) {
	// Valid input
	id, err := GenerateClusterID("test", "infra", "acc")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if id != "test-infra-acc" {
		t.Errorf("expected ID 'test-infra-acc', got %s", id)
	}

	// Missing name
	_, err = GenerateClusterID("", "infra", "acc")
	if err == nil {
		t.Error("expected error for empty name, got nil")
	}

	// Missing accountName
	_, err = GenerateClusterID("test", "infra", "")
	if err == nil {
		t.Error("expected error for empty accountName, got nil")
	}
}

// TestNewCluster tests creation of a new Cluster using NewCluster
func TestNewCluster(t *testing.T) {
	name := "testCluster"
	infraID := "ABCDEF0123"
	provider := AWSProvider
	region := "us-east-1"
	accountName := "testAccount"
	consoleLink := "http://console.testCluster.domain"
	owner := "clusteriq"

	expectedCluster := &Cluster{
		Name:                  name,
		InfraID:               infraID,
		Provider:              provider,
		Status:                Running,
		Region:                region,
		AccountName:           accountName,
		ConsoleLink:           consoleLink,
		InstanceCount:         0,
		Owner:                 owner,
		TotalCost:             0.0,
		Last15DaysCost:        0.0,
		LastMonthCost:         0.0,
		CurrentMonthSoFarCost: 0.0,
		Instances:             make([]Instance, 0),
	}

	actualCluster := NewCluster(name, infraID, provider, region, accountName, consoleLink, owner)

	assert.NotNil(t, actualCluster)

	assert.NotZero(t, actualCluster.LastScanTimestamp)
	assert.Zero(t, actualCluster.InstanceCount)

	expectedCluster.ID = actualCluster.ID
	expectedCluster.Age = actualCluster.Age
	expectedCluster.LastScanTimestamp = actualCluster.LastScanTimestamp
	expectedCluster.CreationTimestamp = actualCluster.CreationTimestamp
	assert.Equal(t, expectedCluster, actualCluster)
}

// TestNewCluster_InvalidParams tests creation of a new Cluster using invalid parameters for the InfraID generation
func TestNewCluster_InvalidParams(t *testing.T) {
	cluster := NewCluster("", "i1", AWSProvider, "us-east-1", "acc1", "https://console", "user")
	if cluster != nil {
		t.Errorf("expected nil cluster when name is empty")
	}

	cluster = NewCluster("name", "i1", AWSProvider, "us-east-1", "", "https://console", "user")
	if cluster != nil {
		t.Errorf("expected nil cluster when account name is empty")
	}
}

// TestIsClusterRunning tests Cluster.IsClusterRunning with both running and non-running status
func TestIsClusterRunning(t *testing.T) {
	// Should return true
	c := Cluster{Status: Running}
	if !c.IsClusterRunning() {
		t.Error("expected IsClusterRunning to return true")
	}

	// Should return false
	c.Status = Stopped
	if c.IsClusterRunning() {
		t.Error("expected IsClusterRunning to return false")
	}
}

// TestIsClusterStopped tests Cluster.IsClusterStopped with both stopped and non-stopped status
func TestIsClusterStopped(t *testing.T) {
	// Should return true
	c := Cluster{Status: Stopped}
	if !c.IsClusterStopped() {
		t.Error("expected IsClusterStopped to return true")
	}

	// Should return false
	c.Status = Running
	if c.IsClusterStopped() {
		t.Error("expected IsClusterStopped to return false")
	}
}

// TestUpdateStatus tests the Cluster.UpdateStatus logic under different scenarios
func TestUpdateStatus(t *testing.T) {
	// Case 1: No instances -> Terminated
	c := Cluster{Instances: []Instance{}}
	c.UpdateStatus()
	if c.Status != Terminated {
		t.Errorf("expected status Terminated, got %v", c.Status)
	}

	// Case 2: At least one Running -> Running
	c.Instances = []Instance{
		{Status: Running},
		{Status: Stopped},
	}
	c.UpdateStatus()
	if c.Status != Running {
		t.Errorf("expected status Running, got %v", c.Status)
	}

	// Case 3: All Terminated -> Terminated
	c.Instances = []Instance{
		{Status: Terminated},
		{Status: Terminated},
	}
	c.UpdateStatus()
	if c.Status != Terminated {
		t.Errorf("expected status Terminated, got %v", c.Status)
	}

	// Case 4: Mix of Stopped and Terminated -> Stopped
	c.Instances = []Instance{
		{Status: Stopped},
		{Status: Terminated},
	}
	c.UpdateStatus()
	if c.Status != Stopped {
		t.Errorf("expected status Stopped, got %v", c.Status)
	}
}

// TestUpdateAge tests Cluster.UpdateAge under valid scenario
func TestUpdateAge(t *testing.T) {
	now := time.Now()
	old := now.Add(-72 * time.Hour)

	c := Cluster{
		LastScanTimestamp: now,
		Age:               3,
		Instances: []Instance{
			{CreationTimestamp: old},
		},
	}
	err := c.UpdateAge()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Age != 3 {
		t.Errorf("expected age 3, got %d", c.Age)
	}
}

// TestUpdateAge_Error tests UpdateAge with decreasing age scenario (should fail)
func TestUpdateAge_Error(t *testing.T) {
	now := time.Now()
	old := now.Add(-24 * time.Hour)

	c := Cluster{
		LastScanTimestamp: now,
		Age:               5,
		Instances: []Instance{
			{CreationTimestamp: old},
		},
	}
	err := c.UpdateAge()
	if err == nil {
		t.Error("expected error due to lower age, got nil")
	}
	if !strings.Contains(err.Error(), "New cluster age is lower") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestClusterUpdateCosts tests Cluster.UpdateCosts including cost validation
func TestClusterUpdateCosts(t *testing.T) {
	// Case 1: total cost too high (should fail)
	c := Cluster{
		TotalCost: 10.0,
		Instances: []Instance{
			{TotalCost: 3.0},
			{TotalCost: 2.0},
		},
	}
	err := c.UpdateCosts()
	if err == nil {
		t.Error("expected error due to lower total cost, got nil", err)
	}

	// Case 2: correct new cost (should succeed)
	c.TotalCost = 5.0
	err = c.UpdateCosts()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.TotalCost != 5.0 {
		t.Errorf("expected total cost 5.0, got %f", c.TotalCost)
	}
}

// TestUpdate tests Cluster.Update end-to-end (status, age, cost)
func TestUpdate(t *testing.T) {
	now := time.Now()
	expectedCost := 5.0

	c := Cluster{
		TotalCost:         expectedCost,
		LastScanTimestamp: now,
		Instances: []Instance{
			{Status: Stopped, CreationTimestamp: now.Add(-72 * time.Hour), TotalCost: 2.5},
			{Status: Terminated, CreationTimestamp: now.Add(-24 * time.Hour), TotalCost: 2.5},
		},
	}

	err := c.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Status != Stopped {
		t.Errorf("expected status Stopped, got %v", c.Status)
	}
	if c.TotalCost != expectedCost {
		t.Errorf("expected total cost %f, got %f", expectedCost, c.TotalCost)
	}
}

// TestUpdate_NoInstances verifies that Update correctly processes the case
// where a cluster doesn't have errors or instances either
func TestUpdate_NoInstances(t *testing.T) {
	now := time.Now()
	expectedCost := 0.0
	c := Cluster{
		Age:               0,
		LastScanTimestamp: now,
		Instances:         []Instance{},
		TotalCost:         expectedCost,
	}

	err := c.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Age != 1 {
		t.Errorf("expected age 0, got %d", c.Age)
	}
	if c.TotalCost != expectedCost {
		t.Errorf("expected total cost %f, got %f", expectedCost, c.TotalCost)
	}
}

// TestUpdate_AgeZero verifies that Update correctly processes the case
// where current age is 0 and the newly calculated age is also 0.
func TestUpdate_AgeZero(t *testing.T) {
	now := time.Now()

	c := Cluster{
		Age:               0, // critical: this skips the error check even if newAge is 0
		LastScanTimestamp: now,
		Instances: []Instance{
			{Status: Running, CreationTimestamp: now},
		},
		TotalCost: 0.0,
	}

	err := c.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Age != 1 {
		t.Errorf("expected age 0, got %d", c.Age)
	}
}

// TestUpdate_ErrorInAge verifies that Cluster.Update returns an error
// when UpdateAge fails due to a decrease in the calculated age.
func TestUpdate_ErrorInAge(t *testing.T) {
	now := time.Now()

	c := Cluster{
		Age:               5, // current age
		LastScanTimestamp: now,
		Instances: []Instance{
			{Status: Stopped, CreationTimestamp: now.Add(-24 * time.Hour), TotalCost: 2.0},
		},
		TotalCost: 2.0, // valid total cost
	}

	err := c.Update()
	if err == nil {
		t.Fatal("expected error from UpdateAge, got nil")
	}
	if !strings.Contains(err.Error(), "New cluster age is lower") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestUpdate_ErrorInCosts verifies that Cluster.Update returns an error
// when UpdateCosts fails because the estimated cost is higher than the current one.
func TestUpdate_ErrorInCosts(t *testing.T) {
	now := time.Now()

	c := Cluster{
		Age:               1, // valid age
		LastScanTimestamp: now,
		Instances: []Instance{
			{Status: Running, CreationTimestamp: now.Add(-24 * time.Hour), TotalCost: 4.0},
		},
		TotalCost: 10.0, // higher than calculated => should trigger error
	}

	err := c.Update()
	if err == nil {
		t.Fatal("expected error from UpdateCosts, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "New estimated cost is lower") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestAddInstance tests Cluster.AddInstance appending correctly and triggering Update
func TestAddInstance(t *testing.T) {
	c := NewCluster("name", "infra", AWSProvider, "region", "acc", "link", "owner")
	if len(c.Instances) != 0 {
		t.Fatalf("expected empty instances list")
	}

	inst := Instance{
		Status:            Running,
		CreationTimestamp: time.Now().Add(-48 * time.Hour),
		TotalCost:         1.5,
	}
	err := c.AddInstance(inst)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(c.Instances) != 1 {
		t.Errorf("expected 1 instance, got %d", len(c.Instances))
	}
}

// TestPrintCluster tests Cluster.PrintCluster (no panic, logs only)
func TestPrintCluster(t *testing.T) {
	c := NewCluster("name", "infra", AWSProvider, "region", "acc", "link", "owner")
	c.PrintCluster()

	now := time.Now()
	c.Instances = []Instance{
		{Status: Running, CreationTimestamp: now.Add(-24 * time.Hour), TotalCost: 4.0},
	}
	c.PrintCluster()
}
