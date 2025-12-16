package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewInstantAction verifies the InstantAction creation.
func TestNewInstantAction(t *testing.T) {
	t.Run("New InstantAction", func(t *testing.T) { testNewInstantAction_Correct(t) })
}

func testNewInstantAction_Correct(t *testing.T) {
	operation := ActionOperation("START")
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1"},
	}
	status := "Pending"
	enabled := true

	action := NewInstantAction(operation, target, status, enabled)

	// Basic check
	assert.NotNil(t, action)

	// Parameters check
	expectedID := target.AccountID + target.ClusterID + string(operation)
	assert.Equal(t, expectedID, action.ID)
	assert.Equal(t, operation, action.Operation)
	assert.Equal(t, target, action.Target)
	assert.Equal(t, status, action.Status)
	assert.Equal(t, enabled, action.Enabled)
}

// TestGetActionOperation verifies the ActionOperation returned by the getter function.
func TestInstantAction_GetActionOperation(t *testing.T) {
	t.Run("GetActionOperation", func(t *testing.T) { testInstantAction_GetActionOperation(t) })
}

func testInstantAction_GetActionOperation(t *testing.T) {
	operation := ActionOperation("STOP")
	action := InstantAction{
		BaseAction: BaseAction{
			Operation: operation,
		},
	}

	assert.Equal(t, operation, action.GetActionOperation())
}

// TestGetRegion verifies the Region returned by the getter function.
func TestInstantAction_GetRegion(t *testing.T) {
	t.Run("GetRegion", func(t *testing.T) { testInstantAction_GetRegion(t) })
}

func testInstantAction_GetRegion(t *testing.T) {
	action := InstantAction{
		BaseAction: BaseAction{
			Target: ActionTarget{
				Region: "us-east-1",
			},
		},
	}

	assert.Equal(t, "us-east-1", action.GetRegion())
}

// TestGetTarget verifies the Target returned by the getter function.
func TestInstantAction_GetTarget(t *testing.T) {
	t.Run("GetTarget", func(t *testing.T) { testInstantAction_GetTarget(t) })
}

func testInstantAction_GetTarget(t *testing.T) {
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1"},
	}
	action := InstantAction{
		BaseAction: BaseAction{
			Target: target,
		},
	}

	assert.Equal(t, target, action.GetTarget())
}

// TestGetID verifies the ID returned by the getter function.
func TestInstantAction_GetID(t *testing.T) {
	t.Run("GetID", func(t *testing.T) { testInstantAction_GetID(t) })
}

func testInstantAction_GetID(t *testing.T) {
	action := InstantAction{
		BaseAction: BaseAction{
			ID: "id-123",
		},
	}

	assert.Equal(t, "id-123", action.GetID())
}

// TestGetType verifies the ActionType returned by the getter function.
func TestInstantAction_GetType(t *testing.T) {
	t.Run("GetType", func(t *testing.T) { testInstantAction_GetType(t) })
}

func testInstantAction_GetType(t *testing.T) {
	action := InstantAction{}

	assert.Equal(t, InstantActionType, action.GetType())
}
