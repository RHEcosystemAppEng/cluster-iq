package actions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewScheduledAction verifies the ScheduledAction creation.
func TestNewScheduledAction(t *testing.T) {
	t.Run("New ScheduledAction", func(t *testing.T) { testNewScheduledAction_Correct(t) })
}

func testNewScheduledAction_Correct(t *testing.T) {
	operation := ActionOperation("START")
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1", "i-2"},
	}
	status := "Pending"
	enabled := true
	when := time.Now().Add(1 * time.Hour)

	action := NewScheduledAction(operation, target, status, enabled, when)

	// Basic check
	assert.NotNil(t, action)

	// Parameters check
	expectedID := target.AccountID + target.ClusterID + string(operation)
	assert.Equal(t, expectedID, action.ID)
	assert.Equal(t, operation, action.Operation)
	assert.Equal(t, target, action.Target)
	assert.Equal(t, status, action.Status)
	assert.Equal(t, enabled, action.Enabled)
	assert.Equal(t, "scheduled_action", action.Type)
	assert.Equal(t, when, action.When)
}

// TestGetActionOperation verifies the ActionOperation returned by the getter function.
func TestScheduledAction_GetActionOperation(t *testing.T) {
	t.Run("GetActionOperation", func(t *testing.T) { testScheduledAction_GetActionOperation(t) })
}

func testScheduledAction_GetActionOperation(t *testing.T) {
	operation := ActionOperation("STOP")
	action := ScheduledAction{
		BaseAction: BaseAction{
			Operation: operation,
		},
	}

	assert.Equal(t, operation, action.GetActionOperation())
}

// TestGetRegion verifies the Region returned by the getter function.
func TestScheduledAction_GetRegion(t *testing.T) {
	t.Run("GetRegion", func(t *testing.T) { testScheduledAction_GetRegion(t) })
}

func testScheduledAction_GetRegion(t *testing.T) {
	action := ScheduledAction{
		BaseAction: BaseAction{
			Target: ActionTarget{
				Region: "us-east-1",
			},
		},
	}

	assert.Equal(t, "us-east-1", action.GetRegion())
}

// TestGetTarget verifies the Target returned by the getter function.
func TestScheduledAction_GetTarget(t *testing.T) {
	t.Run("GetTarget", func(t *testing.T) { testScheduledAction_GetTarget(t) })
}

func testScheduledAction_GetTarget(t *testing.T) {
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1"},
	}
	action := ScheduledAction{
		BaseAction: BaseAction{
			Target: target,
		},
	}

	assert.Equal(t, target, action.GetTarget())
}

// TestGetID verifies the ID returned by the getter function.
func TestScheduledAction_GetID(t *testing.T) {
	t.Run("GetID", func(t *testing.T) { testScheduledAction_GetID(t) })
}

func testScheduledAction_GetID(t *testing.T) {
	action := ScheduledAction{
		BaseAction: BaseAction{
			ID: "id-123",
		},
	}

	assert.Equal(t, "id-123", action.GetID())
}

// TestGetType verifies the ActionType returned by the getter function.
func TestScheduledAction_GetType(t *testing.T) {
	t.Run("GetType", func(t *testing.T) { testScheduledAction_GetType(t) })
}

func testScheduledAction_GetType(t *testing.T) {
	action := ScheduledAction{}

	assert.Equal(t, ScheduledActionType, action.GetType())
}
