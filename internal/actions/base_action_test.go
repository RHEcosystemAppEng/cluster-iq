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
	status := "Pending"
	enabled := true

	action := NewBaseAction(operation, target, status, enabled)

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
