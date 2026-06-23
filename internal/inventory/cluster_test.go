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

// TestCluster_Update covers Cluster.Update() behavior, including status update,
// age update, and error propagation from UpdateAge().
func TestCluster_Update(t *testing.T) {
	t.Run("success updates status and age", testCluster_Update_SuccessUpdatesStatusAndAge)
	t.Run("error from UpdateAge is returned and status is still updated", testCluster_Update_ErrorFromUpdateAgeIsReturnedAndStatusIsStillUpdated)
}

// testCluster_Update_SuccessUpdatesStatusAndAge validates that Update() updates both status and age.
// It also verifies CreatedAt is set to the oldest instance creation timestamp.
func testCluster_Update_SuccessUpdatesStatusAndAge(t *testing.T) {
	c := testClusterWithInstances(
		Instance{Status: Stopped, CreatedAt: time.Now().Add(-48 * time.Hour)},
		Instance{Status: Terminated, CreatedAt: time.Now().Add(-24 * time.Hour)},
	)

	// status must become Stopped (no Running, not all Terminated).
	// CreatedAt must be the oldest instance timestamp.
	expectedCreatedAt := c.Instances[0].CreatedAt
	expectedAge := calculateAge(expectedCreatedAt, time.Now())

	err := c.Update()
	assert.NoError(t, err)
	assert.Equal(t, Stopped, c.Status)
	assert.True(t, c.CreatedAt.Equal(expectedCreatedAt))
	assert.Equal(t, expectedAge, c.Age)
}

// testCluster_Update_ErrorFromUpdateAgeIsReturnedAndStatusIsStillUpdated validates that Update()
// returns UpdateAge() error, but still updates status beforehand.
func testCluster_Update_ErrorFromUpdateAgeIsReturnedAndStatusIsStillUpdated(t *testing.T) {
	oldest := time.Now().Add(-48 * time.Hour)

	c := testClusterWithInstances(
		Instance{Status: Running, CreatedAt: oldest}, // Forces early Running status
		Instance{Status: Stopped, CreatedAt: time.Now().Add(-24 * time.Hour)},
	)

	// force the "new estimated age is lower" error branch.
	newAge := calculateAge(oldest, time.Now())
	c.Age = newAge + 10

	err := c.Update()
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorNewClusterAge.Error())

	// UpdateStatus runs before UpdateAge, so status must be updated even on error.
	assert.Equal(t, Running, c.Status)

	// UpdateAge sets CreatedAt before comparing ages, so CreatedAt is still updated.
	assert.True(t, c.CreatedAt.Equal(oldest))

	// on error, Age must remain unchanged (the assignment happens after the check).
	assert.Equal(t, newAge+10, c.Age)
}

// TestCluster_UpdateAge covers Cluster.UpdateAge() behavior, including empty clusters,
// selecting the oldest instance timestamp, and the "age regression" error.
func TestCluster_UpdateAge(t *testing.T) {
	t.Run("empty cluster sets CreatedAt to now and Age to 1", testCluster_UpdateAge_EmptyClusterSetsCreatedAtToNowAndAgeTo0)
	t.Run("sets CreatedAt to oldest instance and updates Age", testCluster_UpdateAge_SetsCreatedAtToOldestInstanceAndUpdatesAge)
	t.Run("returns error when current Age is greater than newly estimated Age and Age != 0", testCluster_UpdateAge_ReturnsErrorOnAgeRegression)
}

// testCluster_UpdateAge_EmptyClusterSetsCreatedAtToNowAndAgeTo0 validates UpdateAge() behavior
// when the cluster has no instances.
func testCluster_UpdateAge_EmptyClusterSetsCreatedAtToNowAndAgeTo0(t *testing.T) {
	c := &Cluster{}

	err := c.UpdateAge()
	assert.NoError(t, err)

	// with no instances, CreatedAt starts as time.Now() and Age becomes 0 days.
	assert.Equal(t, 1, c.Age)
	assert.False(t, c.CreatedAt.IsZero())
}

// testCluster_UpdateAge_SetsCreatedAtToOldestInstanceAndUpdatesAge validates that UpdateAge()
// uses the oldest instance creation timestamp and updates cluster age accordingly.
func testCluster_UpdateAge_SetsCreatedAtToOldestInstanceAndUpdatesAge(t *testing.T) {
	oldest := time.Now().Add(-72 * time.Hour)
	newer := time.Now().Add(-24 * time.Hour)

	c := testClusterWithInstances(
		Instance{Status: Stopped, CreatedAt: newer},
		Instance{Status: Stopped, CreatedAt: oldest},
	)

	expectedAge := calculateAge(oldest, time.Now())

	err := c.UpdateAge()
	assert.NoError(t, err)

	// CreatedAt must reflect the oldest node in the cluster.
	assert.True(t, c.CreatedAt.Equal(oldest))

	// Age is recalculated using CreatedAt and the scrape time (time.Now()).
	assert.Equal(t, expectedAge, c.Age)
}

// testCluster_UpdateAge_ReturnsErrorOnAgeRegression validates that UpdateAge() fails when it detects
// a regression in age (newly estimated Age is lower than the current Age), and Age != 0.
func testCluster_UpdateAge_ReturnsErrorOnAgeRegression(t *testing.T) {
	oldest := time.Now().Add(-48 * time.Hour)

	c := testClusterWithInstances(
		Instance{Status: Stopped, CreatedAt: oldest},
	)

	newAge := calculateAge(oldest, time.Now())
	c.Age = newAge + 1 // Force regression

	err := c.UpdateAge()
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrorNewClusterAge.Error())

	// CreatedAt is assigned before the regression check.
	assert.True(t, c.CreatedAt.Equal(oldest))

	// on error, Age must remain the previous value.
	assert.Equal(t, newAge+1, c.Age)
}

// TestCluster_UpdateStatus covers Cluster.UpdateStatus() logic for empty clusters,
// Running early return, all Terminated, and all/mixed Stopped/Terminated.
func TestCluster_UpdateStatus(t *testing.T) {
	t.Run("empty cluster -> Terminated", testCluster_UpdateStatus_EmptyClusterTerminated)
	t.Run("any instance Running -> Running (early return)", testCluster_UpdateStatus_AnyRunningEarlyReturnRunning)
	t.Run("all instances Terminated -> Terminated", testCluster_UpdateStatus_AllTerminatedTerminated)
	t.Run("mix of Stopped/Terminated (no Running) -> Stopped", testCluster_UpdateStatus_MixStoppedTerminatedStopped)
	t.Run("all instances Stopped -> Stopped", testCluster_UpdateStatus_AllStoppedStopped)
}

// testCluster_UpdateStatus_EmptyClusterTerminated validates that an empty cluster is marked as Terminated.
func testCluster_UpdateStatus_EmptyClusterTerminated(t *testing.T) {
	c := &Cluster{Instances: nil}

	c.UpdateStatus()
	assert.Equal(t, Terminated, c.Status)
}

// testCluster_UpdateStatus_AnyRunningEarlyReturnRunning validates the early-return rule:
// any Running instance sets the cluster status to Running.
func testCluster_UpdateStatus_AnyRunningEarlyReturnRunning(t *testing.T) {
	c := testClusterWithInstances(
		Instance{Status: Terminated},
		Instance{Status: Running}, // Must win
		Instance{Status: Stopped},
	)

	c.UpdateStatus()
	assert.Equal(t, Running, c.Status)
}

// testCluster_UpdateStatus_AllTerminatedTerminated validates that all Terminated instances
// result in cluster status Terminated.
func testCluster_UpdateStatus_AllTerminatedTerminated(t *testing.T) {
	c := testClusterWithInstances(
		Instance{Status: Terminated},
		Instance{Status: Terminated},
	)

	c.UpdateStatus()
	assert.Equal(t, Terminated, c.Status)
}

// testCluster_UpdateStatus_MixStoppedTerminatedStopped validates that a mix of Stopped/Terminated
// (with no Running instances) results in cluster status Stopped.
func testCluster_UpdateStatus_MixStoppedTerminatedStopped(t *testing.T) {
	c := testClusterWithInstances(
		Instance{Status: Terminated},
		Instance{Status: Stopped},
		Instance{Status: Terminated},
	)

	c.UpdateStatus()
	assert.Equal(t, Stopped, c.Status)
}

// testCluster_UpdateStatus_AllStoppedStopped validates that all Stopped instances
// result in cluster status Stopped.
func testCluster_UpdateStatus_AllStoppedStopped(t *testing.T) {
	c := testClusterWithInstances(
		Instance{Status: Stopped},
		Instance{Status: Stopped},
	)

	c.UpdateStatus()
	assert.Equal(t, Stopped, c.Status)
}

// testClusterWithInstances builds a Cluster with the provided instances.
// It keeps tests concise and consistent across scenarios.
func testClusterWithInstances(instances ...Instance) *Cluster {
	c := &Cluster{
		Instances: make([]Instance, 0, len(instances)),
	}
	c.Instances = append(c.Instances, instances...)
	return c
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
	assert.Equal(t, 0, cluster.InstancesCount())
}

// TestGenerateClusterID tests GenerateClusterID function for 100% coverage
func TestGenerateClusterID(t *testing.T) {
	assert.Equal(t, "test-infra", GenerateClusterID("test", "infra"))
}

