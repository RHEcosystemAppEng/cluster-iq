package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewPowerOnClusterAction verifies NewPowerOnClusterAction builds a valid InstantAction.
func TestNewPowerOnClusterAction(t *testing.T) {
	t.Run("New PowerOnCluster InstantAction", func(t *testing.T) { testNewPowerOnClusterAction_Correct(t) })
}

func testNewPowerOnClusterAction_Correct(t *testing.T) {
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1", "i-2"},
	}
	requester := "user"
	desc := "desc"

	a := NewPowerOnClusterAction(target, requester, &desc)

	assert.NotNil(t, a)
	assert.Equal(t, PowerOnCluster, a.GetActionOperation())
	assert.Equal(t, target, a.GetTarget())
	assert.Equal(t, target.Region, a.GetRegion())
	assert.Equal(t, StatusPending, a.Status)
	assert.True(t, a.Enabled)
	assert.Equal(t, InstantActionType, a.GetType())
	assert.Empty(t, a.GetID())
}

// TestNewPowerOffClusterAction verifies NewPowerOffClusterAction builds a valid InstantAction.
func TestNewPowerOffClusterAction(t *testing.T) {
	t.Run("New PowerOffCluster InstantAction", func(t *testing.T) { testNewPowerOffClusterAction_Correct(t) })
}

func testNewPowerOffClusterAction_Correct(t *testing.T) {
	target := ActionTarget{
		AccountID: "acc-2",
		Region:    "us-east-1",
		ClusterID: "cluster-2",
		Instances: []string{"i-9"},
	}
	requester := "api"
	desc := "power off now"

	a := NewPowerOffClusterAction(target, requester, &desc)

	assert.NotNil(t, a)
	assert.Equal(t, PowerOffCluster, a.GetActionOperation())
	assert.Equal(t, target, a.GetTarget())
	assert.Equal(t, target.Region, a.GetRegion())
	assert.Equal(t, StatusPending, a.Status)
	assert.True(t, a.Enabled)
	assert.Equal(t, InstantActionType, a.GetType())
	assert.Empty(t, a.GetID())
}
