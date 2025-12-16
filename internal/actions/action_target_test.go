package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewActionTarget verifies the ActionTarget creation.
func TestNewActionTarget(t *testing.T) {
	t.Run("New ActionTarget", func(t *testing.T) { testNewActionTarget_Correct(t) })
}

func testNewActionTarget_Correct(t *testing.T) {
	accountID := "acc-1"
	region := "eu-west-1"
	clusterID := "cluster-123"
	instances := []string{"i-1", "i-2"}

	at := NewActionTarget(accountID, region, clusterID, instances)

	// Basic check
	assert.NotNil(t, at)

	// Parameters check
	assert.Equal(t, accountID, at.AccountID)
	assert.Equal(t, region, at.Region)
	assert.Equal(t, clusterID, at.ClusterID)
	assert.Equal(t, instances, at.Instances)
}

// TestGetAccountID verifies the AccountID returned by the getter function.
func TestGetAccountID(t *testing.T) {
	t.Run("GetAccountID", func(t *testing.T) { testGetAccountID(t) })
}

func testGetAccountID(t *testing.T) {
	at := ActionTarget{
		AccountID: "acc-1",
	}

	assert.Equal(t, at.AccountID, at.GetAccountID())
}

// TestGetRegion verifies the Region returned by the getter function.
func TestGetRegion(t *testing.T) {
	t.Run("GetRegion", func(t *testing.T) { testGetRegion(t) })
}

func testGetRegion(t *testing.T) {
	at := ActionTarget{
		Region: "eu-west-1",
	}

	assert.Equal(t, at.Region, at.GetRegion())
}

// TestGetClusterID verifies the ClusterID returned by the getter function.
func TestGetClusterID(t *testing.T) {
	t.Run("GetClusterID", func(t *testing.T) { testGetClusterID(t) })
}

func testGetClusterID(t *testing.T) {
	at := ActionTarget{
		ClusterID: "cluster-123",
	}

	assert.Equal(t, at.ClusterID, at.GetClusterID())
}

// TestGetInstances verifies the Instances returned by the getter function.
func TestGetInstances(t *testing.T) {
	t.Run("GetInstances", func(t *testing.T) { testGetInstances(t) })
}

func testGetInstances(t *testing.T) {
	instances := []string{"i-1", "i-2"}
	at := ActionTarget{
		Instances: instances,
	}

	assert.Equal(t, at.Instances, at.GetInstances())
}
