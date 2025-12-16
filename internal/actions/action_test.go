package actions

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDecodeActions verifies DecodeActions unmarshals actions by type and handles errors.
func TestDecodeActions(t *testing.T) {
	t.Run("Decode InstantAction", func(t *testing.T) { testDecodeActions_InstantAction(t) })
	t.Run("Decode ScheduledAction", func(t *testing.T) { testDecodeActions_ScheduledAction(t) })
	t.Run("Decode CronAction", func(t *testing.T) { testDecodeActions_CronAction(t) })
	t.Run("Decode Mixed actions", func(t *testing.T) { testDecodeActions_Mixed(t) })
	t.Run("Decode empty slice", func(t *testing.T) { testDecodeActions_EmptySlice(t) })
	t.Run("Decode invalid JSON", func(t *testing.T) { testDecodeActions_InvalidJSON(t) })
	t.Run("Decode missing type", func(t *testing.T) { testDecodeActions_MissingType(t) })
	t.Run("Decode unknown type", func(t *testing.T) { testDecodeActions_UnknownType(t) })
	t.Run("Decode InstantAction invalid payload", func(t *testing.T) { testDecodeActions_InstantAction_InvalidPayload(t) })
	t.Run("Decode ScheduledAction invalid payload", func(t *testing.T) { testDecodeActions_ScheduledAction_InvalidPayload(t) })
	t.Run("Decode CronAction invalid payload", func(t *testing.T) { testDecodeActions_CronAction_InvalidPayload(t) })

}
func testDecodeActions_InstantAction(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + InstantActionType + "\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
	assert.Len(t, *actions, 1)

	_, ok := (*actions)[0].(InstantAction)
	assert.True(t, ok)

	raw = []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + InstantActionType + "\", \"time\":\"abcd\"}")),
	}

	actions, err = DecodeActions(raw)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
	assert.Len(t, *actions, 1)

	_, ok = (*actions)[0].(InstantAction)
	assert.True(t, ok)
}

func testDecodeActions_ScheduledAction(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + ScheduledActionType + "\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
	assert.Len(t, *actions, 1)

	_, ok := (*actions)[0].(ScheduledAction)
	assert.True(t, ok)

	raw = []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + ScheduledActionType + "\", \"time\":\"abcd\"}")),
	}

	actions, err = DecodeActions(raw)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
	assert.Len(t, *actions, 1)

	_, ok = (*actions)[0].(ScheduledAction)
	assert.True(t, ok)
}

func testDecodeActions_CronAction(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + CronActionType + "\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
	assert.Len(t, *actions, 1)

	_, ok := (*actions)[0].(CronAction)
	assert.True(t, ok)
}

func testDecodeActions_Mixed(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + ScheduledActionType + "\"}")),
		json.RawMessage([]byte("{\"type\":\"" + CronActionType + "\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
	assert.Len(t, *actions, 2)

	_, ok0 := (*actions)[0].(ScheduledAction)
	_, ok1 := (*actions)[1].(CronAction)
	assert.True(t, ok0)
	assert.True(t, ok1)
}

func testDecodeActions_EmptySlice(t *testing.T) {
	raw := []json.RawMessage{}

	actions, err := DecodeActions(raw)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
	assert.Len(t, *actions, 0)
}

func testDecodeActions_InvalidJSON(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{")),
	}

	actions, err := DecodeActions(raw)
	assert.Error(t, err)
	assert.Nil(t, actions)
}

func testDecodeActions_MissingType(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"foo\":\"bar\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.Error(t, err)
	assert.Nil(t, actions)
	assert.ErrorContains(t, err, "unknown ActionType:")
}

func testDecodeActions_UnknownType(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"UNKNOWN\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.Error(t, err)
	assert.Nil(t, actions)
	assert.ErrorContains(t, err, ErrorUnknownActionType.Error())
}

// TestSplitActionsByType verifies SplitActionsByType splits actions by concrete pointer type.
func TestSplitActionsByType(t *testing.T) {
	t.Run("Split pointers", func(t *testing.T) { testSplitActionsByType_Pointers(t) })
	t.Run("Split values", func(t *testing.T) { testSplitActionsByType_Values(t) })
	t.Run("Split mixed pointer and value", func(t *testing.T) { testSplitActionsByType_Mixed(t) })
}

func testSplitActionsByType_Pointers(t *testing.T) {
	s1 := &ScheduledAction{}
	c1 := &CronAction{}

	actions := []Action{s1, c1}

	sched, cron := SplitActionsByType(actions)

	assert.Len(t, sched, 1)
	assert.Len(t, cron, 1)
	assert.Same(t, s1, sched[0])
	assert.Same(t, c1, cron[0])
}

func testSplitActionsByType_Values(t *testing.T) {
	// SplitActionsByType matches *ScheduledAction and *CronAction only.
	actions := []Action{
		ScheduledAction{},
		CronAction{},
	}

	sched, cron := SplitActionsByType(actions)

	assert.Len(t, sched, 0)
	assert.Len(t, cron, 0)
}

func testSplitActionsByType_Mixed(t *testing.T) {
	s1 := &ScheduledAction{}

	actions := []Action{
		s1,
		CronAction{}, // value should be ignored by the type switch
	}

	sched, cron := SplitActionsByType(actions)

	assert.Len(t, sched, 1)
	assert.Len(t, cron, 0)
	assert.Same(t, s1, sched[0])
}

func testDecodeActions_InstantAction_InvalidPayload(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + InstantActionType + "\",\"enabled\":\"nope\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.Error(t, err)
	assert.Nil(t, actions)
}

func testDecodeActions_ScheduledAction_InvalidPayload(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + ScheduledActionType + "\",\"enabled\":\"nope\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.Error(t, err)
	assert.Nil(t, actions)
}

func testDecodeActions_CronAction_InvalidPayload(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage([]byte("{\"type\":\"" + CronActionType + "\",\"enabled\":\"nope\"}")),
	}

	actions, err := DecodeActions(raw)
	assert.Error(t, err)
	assert.Nil(t, actions)
}
