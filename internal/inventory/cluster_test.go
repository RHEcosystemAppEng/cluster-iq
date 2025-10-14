package inventory

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: Group by function and test
// TODO: Include Asserts
// TODO: Comments

// TestGenerateClusterID tests GenerateClusterID with valid and invalid inputs
func TestGenerateClusterID(t *testing.T) {
	// Valid input
	assert.Equal(t, "test-infra", generateClusterID("test", "infra"))
}

// TestNewCluster tests creation of a new Cluster using NewCluster
func TestNewCluster(t *testing.T) {
	name := "testCluster"
	infraID := "ABCDEF0123"
	provider := AWSProvider
	region := "us-east-1"
	consoleLink := "http://console.testCluster.domain"
	owner := "clusteriq"

	expectedCluster := &Cluster{
		ClusterName: name,
		InfraID:     infraID,
		Provider:    provider,
		Status:      Running,
		Region:      region,
		AccountID:   "",
		ConsoleLink: consoleLink,
		Owner:       owner,
		Instances:   make([]Instance, 0),
	}

	actualCluster := NewCluster(name, infraID, provider, region, consoleLink, owner)

	assert.NotNil(t, actualCluster)

	assert.Zero(t, actualCluster.LastScanTS)

	expectedCluster.ClusterID = actualCluster.ClusterID
	expectedCluster.Age = actualCluster.Age
	expectedCluster.LastScanTS = actualCluster.LastScanTS
	expectedCluster.CreatedAt = actualCluster.CreatedAt
	assert.Equal(t, expectedCluster, actualCluster)
}

// TestNewCluster_InvalidParams tests creation of a new Cluster using invalid parameters for the InfraID generation
func TestNewCluster_InvalidParams(t *testing.T) {
	cluster := NewCluster("", "i1", AWSProvider, "us-east-1", "https://console", "user")
	if cluster != nil {
		t.Errorf("expected nil cluster when name is empty")
	}
}

// TestIsClusterRunning tests Cluster.IsClusterRunning with both running and non-running status
func TestIsClusterRunning(t *testing.T) {
	// Should return true
	cluster := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.True(t, cluster.IsClusterRunning())

	// Should return false
	cluster.Status = Stopped
	assert.False(t, cluster.IsClusterRunning())
}

// TestIsClusterStopped tests Cluster.IsClusterStopped with both stopped and non-stopped status
func TestIsClusterStopped(t *testing.T) {
	// Should return true
	cluster := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.False(t, cluster.IsClusterStopped())

	// Should return false
	cluster.Status = Stopped
	assert.True(t, cluster.IsClusterStopped())
}

// TestUpdateStatus tests the Cluster.UpdateStatus logic under different scenarios
func TestUpdateStatus(t *testing.T) {
	// Case 1: No instances -> Terminated
	cluster := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	cluster.UpdateStatus()
	if cluster.Status != Terminated {
		t.Errorf("expected status Terminated, got %v", cluster.Status)
	}

	// Case 2: At least one Running -> Running
	cluster.Instances = []Instance{
		{Status: Running},
		{Status: Stopped},
	}
	cluster.UpdateStatus()
	if cluster.Status != Running {
		t.Errorf("expected status Running, got %v", cluster.Status)
	}

	// Case 3: All Terminated -> Terminated
	cluster.Instances = []Instance{
		{Status: Terminated},
		{Status: Terminated},
	}
	cluster.UpdateStatus()
	if cluster.Status != Terminated {
		t.Errorf("expected status Terminated, got %v", cluster.Status)
	}

	// Case 4: Mix of Stopped and Terminated -> Stopped
	cluster.Instances = []Instance{
		{Status: Stopped},
		{Status: Terminated},
	}
	cluster.UpdateStatus()
	if cluster.Status != Stopped {
		t.Errorf("expected status Stopped, got %v", cluster.Status)
	}
}

// TestUpdateAge tests Cluster.UpdateAge under valid scenario
func TestUpdateAge(t *testing.T) {
	now := time.Now()
	old := now.Add(-72 * time.Hour)

	c := Cluster{
		LastScanTS: now,
		Age:        3,
		Instances: []Instance{
			{CreatedAt: old},
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
		LastScanTS: now,
		Age:        5,
		Instances: []Instance{
			{CreatedAt: old},
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

// TestUpdate tests Cluster.Update end-to-end (status, age, cost)
func TestUpdate(t *testing.T) {
	now := time.Now()

	c := Cluster{
		LastScanTS: now,
		Instances: []Instance{
			{Status: Stopped, CreatedAt: now.Add(-72 * time.Hour)},
			{Status: Terminated, CreatedAt: now.Add(-24 * time.Hour)},
		},
	}

	err := c.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Status != Stopped {
		t.Errorf("expected status Stopped, got %v", c.Status)
	}
}

// TestUpdate_NoInstances verifies that Update correctly processes the case
// where a cluster doesn't have errors or instances either
func TestUpdate_NoInstances(t *testing.T) {
	now := time.Now()
	c := Cluster{
		Age:        0,
		LastScanTS: now,
		Instances:  []Instance{},
	}

	err := c.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Age != 1 {
		t.Errorf("expected age 0, got %d", c.Age)
	}
}

// TestUpdate_AgeZero verifies that Update correctly processes the case
// where current age is 0 and the newly calculated age is also 0.
func TestUpdate_AgeZero(t *testing.T) {
	now := time.Now()

	c := Cluster{
		Age:        0, // critical: this skips the error check even if newAge is 0
		LastScanTS: now,
		Instances: []Instance{
			{Status: Running, CreatedAt: now},
		},
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
		Age:        5, // current age
		LastScanTS: now,
		Instances: []Instance{
			{Status: Stopped, CreatedAt: now.Add(-24 * time.Hour)},
		},
	}

	err := c.Update()
	if err == nil {
		t.Fatal("expected error from UpdateAge, got nil")
	}
	if !strings.Contains(err.Error(), "New cluster age is lower") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestAddInstance tests Cluster.AddInstance appending correctly and triggering Update
func TestAddInstance(t *testing.T) {
	c := NewCluster("name", "infra", AWSProvider, "region", "link", "owner")
	if len(c.Instances) != 0 {
		t.Fatalf("expected empty instances list")
	}

	inst := Instance{
		Status:    Running,
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}
	err := c.AddInstance(&inst)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(c.Instances) != 1 {
		t.Errorf("expected 1 instance, got %d", len(c.Instances))
	}
}

// TestPrintCluster tests Cluster.PrintCluster (no panic, logs only)
func TestPrintCluster(t *testing.T) {
	c := NewCluster("name", "infra", AWSProvider, "region", "link", "owner")
	c.PrintCluster()

	now := time.Now()
	c.Instances = []Instance{
		{Status: Running, CreatedAt: now.Add(-24 * time.Hour)},
	}
	c.PrintCluster()
}
