package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewBaseAction verifies the BaseAction creation.
func TestNewBaseAction(t *testing.T) {
	t.Run("New BaseAction", func(t *testing.T) { testNewBaseAction_Correct(t) })
}

func testNewBaseAction_Correct(t *testing.T) {
	operation := ActionOperation("START")
	target := ActionTarget{
		AccountID: "acc-1",
		Region:    "eu-west-1",
		ClusterID: "cluster-1",
		Instances: []string{"i-1"},
	}
	requester := "someone"
	description := "description"
	status := StatusRunning
	enabled := true

	action := NewBaseAction(operation, target, status, requester, &description, enabled)

	// Basic check
	assert.NotNil(t, action)

	// Parameters check
	assert.Equal(t, operation, action.Operation)
	assert.Equal(t, target, action.Target)
	assert.Equal(t, status, action.Status)
	assert.Equal(t, enabled, action.Enabled)
}

// TestGetActionOperation verifies the ActionOperation returned by the getter function.
func TestGetActionOperation(t *testing.T) {
	t.Run("GetActionOperation", func(t *testing.T) { testGetActionOperation(t) })
}

func testGetActionOperation(t *testing.T) {
	operation := ActionOperation("STOP")
	action := BaseAction{
		Operation: operation,
	}

	assert.Equal(t, operation, action.GetActionOperation())
}

// TestBaseAction_SetStatus verifies that SetStatus updates the action status.
func TestBaseAction_SetStatus(t *testing.T) {
	t.Run("SetStatus updates status correctly", func(t *testing.T) {
		testBaseAction_SetStatus_Correct(t)
	})
}

func testBaseAction_SetStatus_Correct(t *testing.T) {
	action := BaseAction{
		Status: StatusPending,
	}

	assert.Equal(t, StatusPending, action.Status)

	action.SetStatus(StatusCompleted)

	assert.Equal(t, StatusCompleted, action.Status)
}
