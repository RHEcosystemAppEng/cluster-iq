package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCronAction verifies the CronAction creation.
func TestNewCronAction(t *testing.T) {
	t.Run("New CronAction", func(t *testing.T) { testNewCronAction_Correct(t) })
}

func testNewCronAction_Correct(t *testing.T) {
	operation := ActionOperation("START")
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1", "i-2"},
	}
	status := "Pending"
	enabled := true
	expression := "0 0 * * *"

	action := NewCronAction(operation, target, status, enabled, expression)

	// Basic check
	assert.NotNil(t, action)

	// Parameters check
	expectedID := target.AccountID + target.ClusterID + string(operation)
	assert.Equal(t, expectedID, action.ID)
	assert.Equal(t, operation, action.Operation)
	assert.Equal(t, target, action.Target)
	assert.Equal(t, status, action.Status)
	assert.Equal(t, enabled, action.Enabled)
	assert.Equal(t, "cron_action", action.Type)
	assert.Equal(t, expression, action.Expression)
}

// TestGetActionOperation verifies the ActionOperation returned by the getter function.
func TestCronAction_GetActionOperation(t *testing.T) {
	t.Run("GetActionOperation", func(t *testing.T) { testCronAction_GetActionOperation(t) })
}

func testCronAction_GetActionOperation(t *testing.T) {
	operation := ActionOperation("STOP")
	action := CronAction{
		BaseAction: BaseAction{
			Operation: operation,
		},
	}

	assert.Equal(t, operation, action.GetActionOperation())
}

// TestGetRegion verifies the Region returned by the getter function.
func TestCronAction_GetRegion(t *testing.T) {
	t.Run("GetRegion", func(t *testing.T) { testCronAction_GetRegion(t) })
}

func testCronAction_GetRegion(t *testing.T) {
	action := CronAction{
		BaseAction: BaseAction{
			Target: ActionTarget{
				Region: "us-east-1",
			},
		},
	}

	assert.Equal(t, "us-east-1", action.GetRegion())
}

// TestGetTarget verifies the Target returned by the getter function.
func TestCronAction_GetTarget(t *testing.T) {
	t.Run("GetTarget", func(t *testing.T) { testCronAction_GetTarget(t) })
}

func testCronAction_GetTarget(t *testing.T) {
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1"},
	}
	action := CronAction{
		BaseAction: BaseAction{
			Target: target,
		},
	}

	assert.Equal(t, target, action.GetTarget())
}

// TestGetID verifies the ID returned by the getter function.
func TestCronAction_GetID(t *testing.T) {
	t.Run("GetID", func(t *testing.T) { testCronAction_GetID(t) })
}

func testCronAction_GetID(t *testing.T) {
	action := CronAction{
		BaseAction: BaseAction{
			ID: "id-123",
		},
	}

	assert.Equal(t, "id-123", action.GetID())
}

// TestGetType verifies the ActionType returned by the getter function.
func TestCronAction_GetType(t *testing.T) {
	t.Run("GetType", func(t *testing.T) { testCronAction_GetType(t) })
}

func testCronAction_GetType(t *testing.T) {
	action := CronAction{}

	assert.Equal(t, CronActionType, action.GetType())
}

// TestGetCronExpression verifies the cron expression returned by the getter function.
func TestCronAction_GetCronExpression(t *testing.T) {
	t.Run("GetCronExpression", func(t *testing.T) { testCronAction_GetCronExpression(t) })
}

func testCronAction_GetCronExpression(t *testing.T) {
	action := CronAction{
		Expression: "*/5 * * * *",
	}

	assert.Equal(t, "*/5 * * * *", action.GetCronExpression())
}
