package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsResourceStatus(t *testing.T) {
	t.Run("Running", func(t *testing.T) { testRunningStatus(t) })
	t.Run("Stopped", func(t *testing.T) { testStoppedStatus(t) })
	t.Run("Stop", func(t *testing.T) { testStopStatus(t) })
	t.Run("Terminated", func(t *testing.T) { testTerminatedStatus(t) })
	t.Run("Anything", func(t *testing.T) { testAnythingStatus(t) })
}

func testRunningStatus(t *testing.T) {
	assert.Equal(t, AsResourceStatus("RUNNING"), Running)
	assert.Equal(t, AsResourceStatus("Running"), Running)
	assert.Equal(t, AsResourceStatus("running"), Running)
}

func testStoppedStatus(t *testing.T) {
	assert.Equal(t, AsResourceStatus("STOPPED"), Stopped)
	assert.Equal(t, AsResourceStatus("Stopped"), Stopped)
	assert.Equal(t, AsResourceStatus("stopped"), Stopped)
}

func testStopStatus(t *testing.T) {
	assert.Equal(t, AsResourceStatus("STOP"), Stopped)
	assert.Equal(t, AsResourceStatus("Stop"), Stopped)
	assert.Equal(t, AsResourceStatus("stop"), Stopped)
}

func testTerminatedStatus(t *testing.T) {
	assert.Equal(t, AsResourceStatus("TERMINATED"), Terminated)
	assert.Equal(t, AsResourceStatus("Terminated"), Terminated)
	assert.Equal(t, AsResourceStatus("terminated"), Terminated)
}

func testAnythingStatus(t *testing.T) {
	assert.Equal(t, AsResourceStatus("anything"), Running)
	assert.Equal(t, AsResourceStatus("Unknown"), Running)
}
