package inventory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewCluster tests creation of a new Cluster using NewCluster
func TestNewCluster(t *testing.T) {
	t.Run("New Cluster", func(t *testing.T) { testNewCluster_Correct(t) })
	t.Run("New Cluster without ClusterName", func(t *testing.T) { testNewCluster_WithoutClusterName(t) })
	t.Run("New Cluster without InfraID", func(t *testing.T) { testNewCluster_WithoutInfraID(t) })
}

func testNewCluster_Correct(t *testing.T) {
	name := "testCluster"
	infraID := "ABCDEF0123"
	provider := AWSProvider
	region := "us-east-1"
	consoleLink := "http://console.testCluster.domain"
	owner := "clusteriq"

	cluster, err := NewCluster(name, infraID, provider, region, consoleLink, owner)

	// Basic check
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	// Parameters Check
	assert.Equal(t, cluster.ClusterID, GenerateClusterID(cluster.ClusterName, cluster.InfraID))
	assert.Equal(t, cluster.ClusterName, name)
	assert.Equal(t, cluster.InfraID, infraID)
	assert.Equal(t, cluster.Provider, provider)
	assert.Equal(t, cluster.Status, Stopped)
	assert.Equal(t, cluster.Region, region)
	assert.Equal(t, cluster.AccountID, "")
	assert.Equal(t, cluster.ConsoleLink, consoleLink)
	assert.False(t, cluster.LastScanTimestamp.IsZero())
	assert.False(t, cluster.CreatedAt.IsZero())
	assert.Zero(t, cluster.Age)
	assert.Equal(t, cluster.Owner, owner)
	assert.NotNil(t, cluster.Instances)
}

func testNewCluster_WithoutClusterName(t *testing.T) {
	cluster, err := NewCluster("", "ABCDEF0123", AWSProvider, "us-east-1", "https://console", "owner")

	assert.Nil(t, cluster)
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorMissingClusterNameCreation.Error())
}

func testNewCluster_WithoutInfraID(t *testing.T) {
	cluster, err := NewCluster("testCluster", "", AWSProvider, "us-east-1", "https://console", "owner")

	assert.Nil(t, cluster)
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorMissingClusterInfraIDCreation.Error())
}

// TestIsClusterStopped tests the IsClusterStopped function
func TestIsClusterStopped(t *testing.T) {
	t.Run("Cluster Running", func(t *testing.T) { testIsClusterStopped_Running(t) })
	t.Run("Cluster Stopped", func(t *testing.T) { testIsClusterStopped_Stopped(t) })
	t.Run("Cluster Terminated", func(t *testing.T) { testIsClusterStopped_Terminated(t) })
}

func testIsClusterStopped_Running(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	cluster.Status = Stopped
	assert.True(t, cluster.IsClusterStopped())
}

func testIsClusterStopped_Stopped(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	cluster.Status = Running
	assert.False(t, cluster.IsClusterStopped())
}

func testIsClusterStopped_Terminated(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	cluster.Status = Terminated
	assert.False(t, cluster.IsClusterStopped())
}

// TestIsClusterRunning tests the IsClusterRunning function
func TestIsClusterRunning(t *testing.T) {
	t.Run("Cluster Running", func(t *testing.T) { testIsClusterRunning_Running(t) })
	t.Run("Cluster Stopped", func(t *testing.T) { testIsClusterRunning_Stopped(t) })
	t.Run("Cluster Terminated", func(t *testing.T) { testIsClusterRunning_Terminated(t) })
}

func testIsClusterRunning_Running(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	cluster.Status = Stopped
	assert.False(t, cluster.IsClusterRunning())
}

func testIsClusterRunning_Stopped(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	cluster.Status = Running
	assert.True(t, cluster.IsClusterRunning())
}

func testIsClusterRunning_Terminated(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	cluster.Status = Terminated
	assert.False(t, cluster.IsClusterRunning())
}

// TestUpdate tests the Update function including Age and Status updates
func TestUpdate(t *testing.T) {
	t.Run("UpdateAge", func(t *testing.T) { testUpdate_Correct(t) })
	t.Run("UpdateAge Lower New Age", func(t *testing.T) { testUpdate_NewerAge(t) })
}

func testUpdate_Correct(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	err = cluster.Update()
	assert.Nil(t, err)
}

func testUpdate_NewerAge(t *testing.T) {
	prevAge := 30
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)
	t1 := time.Now()
	t2 := t1.Add(-12 * 24 * time.Hour)
	cluster.Instances = []Instance{
		{CreatedAt: t1},
		{CreatedAt: t2},
	}

	cluster.Age = prevAge
	err = cluster.Update()
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorNewClusterAge.Error())
}

// TestUpdateAge tests UpdateAge function
func TestUpdateAge(t *testing.T) {
	t.Run("UpdateAge", func(t *testing.T) { testUpdateAge_Correct(t) })
	t.Run("UpdateAge Lower New Age", func(t *testing.T) { testUpdateAge_NewerAge(t) })
}

func testUpdateAge_Correct(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)
	t1 := time.Now()
	t2 := t1.Add(-12 * 24 * time.Hour)
	cluster.Instances = []Instance{
		{CreatedAt: t1},
		{CreatedAt: t2},
	}

	err = cluster.UpdateAge()
	assert.Nil(t, err)
	assert.Equal(t, cluster.CreatedAt, t2)
	assert.Equal(t, cluster.Age, 11)
}

func testUpdateAge_NewerAge(t *testing.T) {
	prevAge := 30
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)
	t1 := time.Now()
	t2 := t1.Add(-12 * 24 * time.Hour)
	cluster.Instances = []Instance{
		{CreatedAt: t1},
		{CreatedAt: t2},
	}
	cluster.Age = prevAge

	err = cluster.UpdateAge()
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorNewClusterAge.Error())
	assert.NotZero(t, cluster.CreatedAt)
	assert.Equal(t, cluster.Age, prevAge)
}

// TestUpdateStatus tests the UpdateStatus function
func TestUpdateStatus(t *testing.T) {
	t.Run("UpdateStatus no Instances", func(t *testing.T) { testUpdateStatus_NoInstances(t) })
	t.Run("UpdateStatus at lease one Instance Running", func(t *testing.T) { testUpdateStatus_AtLeastOneInstanceRunning(t) })
	t.Run("UpdateStatus all Instances Terminated", func(t *testing.T) { testUpdateStatus_TerminatedInstances(t) })
	t.Run("UpdateStatus mix Stopped and Terminated", func(t *testing.T) { testUpdateStatus_NoInstancesRunning(t) })
}

func testUpdateStatus_NoInstances(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	cluster.UpdateStatus()
	assert.Equal(t, cluster.Status, Terminated)
}

func testUpdateStatus_AtLeastOneInstanceRunning(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)
	cluster.Instances = []Instance{
		{Status: Running},
		{Status: Stopped},
	}

	cluster.UpdateStatus()
	assert.Equal(t, cluster.Status, Running)
}

func testUpdateStatus_TerminatedInstances(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)
	cluster.Instances = []Instance{
		{Status: Terminated},
		{Status: Terminated},
	}

	cluster.UpdateStatus()
	assert.Equal(t, cluster.Status, Terminated)
}

func testUpdateStatus_NoInstancesRunning(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)
	cluster.Instances = []Instance{
		{Status: Stopped},
		{Status: Terminated},
	}

	cluster.UpdateStatus()
	assert.Equal(t, cluster.Status, Stopped)
}

// TestAddInstance tests the AddInstance function including repeated instances
func TestAddInstance(t *testing.T) {
	t.Run("Add Instance", func(t *testing.T) { testAddInstance_Correct(t) })
	t.Run("Add Instance twice", func(t *testing.T) { testAddInstance_Twice(t) })
}

func testAddInstance_Correct(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	instance := Instance{
		InstanceName: "A1",
		ClusterID:    "",
	}

	err = cluster.AddInstance(&instance)
	assert.Nil(t, err)
	assert.Equal(t, instance.ClusterID, cluster.ClusterID)
	assert.Equal(t, len(cluster.Instances), 1)
}

func testAddInstance_Twice(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.Nil(t, err)
	assert.NotNil(t, cluster)

	instance := Instance{
		InstanceName: "A1",
		ClusterID:    "",
	}

	err = cluster.AddInstance(&instance)
	assert.Nil(t, err)
	assert.Equal(t, instance.ClusterID, cluster.ClusterID)
	assert.Equal(t, len(cluster.Instances), 1)

	err = cluster.AddInstance(&instance)
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorAddingInstanceToCluster.Error())
	assert.Equal(t, len(cluster.Instances), 1)
}

// TestDeleteInstance tests the DeleteInstance function including missing instances
func TestDeleteInstance(t *testing.T) {
	t.Run("Delete Instance", func(t *testing.T) { testDeleteInstance_Correct(t) })
	t.Run("Delete missing Instance", func(t *testing.T) { testDeleteInstance_MissingInstance(t) })
}

func testDeleteInstance_Correct(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	instances := []Instance{
		{InstanceName: "A1", ClusterID: "testCluster-i1"},
	}
	cluster.Instances = instances

	err = cluster.DeleteInstance(instances[0].InstanceName)
	assert.Nil(t, err)
	assert.Zero(t, len(cluster.Instances))
}

func testDeleteInstance_MissingInstance(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)

	err = cluster.DeleteInstance("testInstance")
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorDeleteInstanceFromCluster.Error())
	assert.Zero(t, len(cluster.Instances))
}

// TestInstancesCount tests the InstanceCount function for 100% coverage
func TestInstancesCount(t *testing.T) {
	t.Run("Instances Count", func(t *testing.T) { testInstancesCount(t) })
}

func testInstancesCount(t *testing.T) {
	cluster, err := NewCluster("testCluster", "i1", AWSProvider, "us-east-1", "https://console", "user")
	assert.NotNil(t, cluster)
	assert.Nil(t, err)
	assert.Equal(t, len(cluster.Instances), 0)
}

// TestGenerateClusterID tests GenerateClusterID function for 100% coverage
func TestGenerateClusterID(t *testing.T) {
	assert.Equal(t, "test-infra", GenerateClusterID("test", "infra"))
}

// TestPrintCluster tests PrintCluster function for 100% coverage
func TestPrintCluster(t *testing.T) {
	c, _ := NewCluster("name", "infra", AWSProvider, "region", "link", "owner")
	c.PrintCluster()

	now := time.Now()
	c.Instances = []Instance{
		{Status: Running, CreatedAt: now.Add(-24 * time.Hour)},
	}
	c.PrintCluster()
}
