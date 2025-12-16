package dto

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/stretchr/testify/assert"
)

// TestEventDTORequest_ToModelEvent verifies DTO to events.Event conversion.
func TestEventDTORequest_ToModelEvent(t *testing.T) {
	t.Run("Convert DTO to model Event", func(t *testing.T) { testEventDTORequest_ToModelEvent_Correct(t) })
}

func testEventDTORequest_ToModelEvent_Correct(t *testing.T) {
	now := time.Now().UTC()
	desc := "event description"

	dto := EventDTORequest{
		ID:             42,
		Action:         "START",
		ResourceID:     "cluster-1",
		ResourceType:   "cluster",
		EventTimestamp: now,
		Result:         "Success",
		Severity:       "Info",
		TriggeredBy:    "scanner",
		Description:    &desc,
	}

	event := dto.ToModelEvent()

	assert.NotNil(t, event)
	assert.Equal(t, dto.ID, event.ID)
	assert.Equal(t, actions.ActionOperation(dto.Action), event.Action)
	assert.Equal(t, dto.EventTimestamp, event.EventTimestamp)
	assert.Equal(t, dto.Description, event.Description)
	assert.Equal(t, dto.ResourceID, event.ResourceID)
	assert.Equal(t, dto.ResourceType, event.ResourceType)
	assert.Equal(t, dto.Result, event.Result)
	assert.Equal(t, dto.Severity, event.Severity)
	assert.Equal(t, dto.TriggeredBy, event.TriggeredBy)
}

// TestEventDTORequest_ToModelEvent_NilDescription verifies nil description handling.
func TestEventDTORequest_ToModelEvent_NilDescription(t *testing.T) {
	t.Run("Nil description", func(t *testing.T) { testEventDTORequest_ToModelEvent_NilDescription(t) })
}

func testEventDTORequest_ToModelEvent_NilDescription(t *testing.T) {
	now := time.Now().UTC()

	dto := EventDTORequest{
		ID:             7,
		Action:         "STOP",
		ResourceID:     "instance-1",
		ResourceType:   "instance",
		EventTimestamp: now,
		Result:         "Failed",
		Severity:       "Error",
		TriggeredBy:    "agent",
		Description:    nil,
	}

	event := dto.ToModelEvent()

	assert.NotNil(t, event)
	assert.Nil(t, event.Description)
}

// TestEventDTOResponse_Struct verifies DTO response structs can be instantiated.
func TestEventDTOResponse_Struct(t *testing.T) {
	t.Run("ClusterEventDTOResponse struct", func(t *testing.T) { testClusterEventDTOResponse_Struct(t) })
	t.Run("SystemEventDTOResponse struct", func(t *testing.T) { testSystemEventDTOResponse_Struct(t) })
}

func testClusterEventDTOResponse_Struct(t *testing.T) {
	desc := "desc"
	now := time.Now()

	dto := ClusterEventDTOResponse{
		ID:             1,
		Action:         "START",
		ResourceID:     "cluster-1",
		ResourceType:   "cluster",
		EventTimestamp: now,
		Result:         "Success",
		Severity:       "Info",
		TriggeredBy:    "api",
		Description:    &desc,
	}

	assert.Equal(t, int64(1), dto.ID)
	assert.Equal(t, "START", dto.Action)
	assert.Equal(t, "cluster-1", dto.ResourceID)
	assert.Equal(t, "cluster", dto.ResourceType)
	assert.Equal(t, now, dto.EventTimestamp)
	assert.Equal(t, "Success", dto.Result)
	assert.Equal(t, "Info", dto.Severity)
	assert.Equal(t, "api", dto.TriggeredBy)
	assert.Equal(t, &desc, dto.Description)
}

func testSystemEventDTOResponse_Struct(t *testing.T) {
	desc := "desc"
	now := time.Now()

	dto := SystemEventDTOResponse{
		ClusterEventDTOResponse: ClusterEventDTOResponse{
			ID:             2,
			Action:         "STOP",
			ResourceID:     "cluster-2",
			ResourceType:   "cluster",
			EventTimestamp: now,
			Result:         "Failed",
			Severity:       "Error",
			TriggeredBy:    "scheduler",
			Description:    &desc,
		},
		AccountID: "acc-1",
		Provider:  "AWS",
	}

	assert.Equal(t, int64(2), dto.ID)
	assert.Equal(t, "STOP", dto.Action)
	assert.Equal(t, "cluster-2", dto.ResourceID)
	assert.Equal(t, "cluster", dto.ResourceType)
	assert.Equal(t, now, dto.EventTimestamp)
	assert.Equal(t, "Failed", dto.Result)
	assert.Equal(t, "Error", dto.Severity)
	assert.Equal(t, "scheduler", dto.TriggeredBy)
	assert.Equal(t, &desc, dto.Description)
	assert.Equal(t, "acc-1", dto.AccountID)
	assert.Equal(t, "AWS", dto.Provider)
}

// Compile-time check to ensure returned type matches expected model
func TestEventDTORequest_ReturnType(t *testing.T) {
	t.Run("Return type is *events.Event", func(t *testing.T) { testEventDTORequest_ReturnType(t) })
}

func testEventDTORequest_ReturnType(t *testing.T) {
	dto := EventDTORequest{}
	event := dto.ToModelEvent()

	_, ok := interface{}(event).(*events.Event)
	assert.True(t, ok)
}
