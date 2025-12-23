package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestClusterEventDBResponse_ToClusterEventDTOResponse verifies DB model to DTO conversion.
func TestClusterEventDBResponse_ToClusterEventDTOResponse(t *testing.T) {
	t.Run("Convert ClusterEventDBResponse to ClusterEventDTOResponse", func(t *testing.T) {
		testClusterEventDBResponse_ToClusterEventDTOResponse_Correct(t)
	})
	t.Run("Convert with nil description", func(t *testing.T) {
		testClusterEventDBResponse_ToClusterEventDTOResponse_NilDescription(t)
	})
}

func testClusterEventDBResponse_ToClusterEventDTOResponse_Correct(t *testing.T) {
	now := time.Now().UTC()
	desc := "desc"

	model := ClusterEventDBResponse{
		ID:             1,
		EventTimestamp: now,
		TriggeredBy:    "api",
		Action:         "START",
		ResourceID:     "cluster-1",
		ResourceType:   "cluster",
		Result:         "Success",
		Description:    &desc,
		Severity:       "Info",
	}

	dto := model.ToClusterEventDTOResponse()

	assert.NotNil(t, dto)
	assert.Equal(t, model.ID, dto.ID)
	assert.Equal(t, model.EventTimestamp, dto.EventTimestamp)
	assert.Equal(t, model.TriggeredBy, dto.TriggeredBy)
	assert.Equal(t, model.Action, dto.Action)
	assert.Equal(t, model.ResourceID, dto.ResourceID)
	assert.Equal(t, model.ResourceType, dto.ResourceType)
	assert.Equal(t, model.Result, dto.Result)
	assert.Equal(t, model.Description, dto.Description)
	assert.Equal(t, model.Severity, dto.Severity)
}

func testClusterEventDBResponse_ToClusterEventDTOResponse_NilDescription(t *testing.T) {
	now := time.Now().UTC()

	model := ClusterEventDBResponse{
		ID:             2,
		EventTimestamp: now,
		TriggeredBy:    "scanner",
		Action:         "STOP",
		ResourceID:     "cluster-2",
		ResourceType:   "cluster",
		Result:         "Failed",
		Description:    nil,
		Severity:       "Error",
	}

	dto := model.ToClusterEventDTOResponse()

	assert.NotNil(t, dto)
	assert.Equal(t, model.ID, dto.ID)
	assert.Nil(t, dto.Description)
}

// TestToClusterEventDTOResponseList verifies list conversion from DB responses to DTOs.
func TestToClusterEventDTOResponseList(t *testing.T) {
	t.Run("Convert ClusterEventDBResponse list to DTO list", func(t *testing.T) {
		testToClusterEventDTOResponseList_Correct(t)
	})
}

func testToClusterEventDTOResponseList_Correct(t *testing.T) {
	now := time.Now().UTC()

	models := []ClusterEventDBResponse{
		{ID: 1, EventTimestamp: now, TriggeredBy: "api", Action: "START", ResourceID: "c1", ResourceType: "cluster", Result: "Success", Severity: "Info"},
		{ID: 2, EventTimestamp: now.Add(-time.Minute), TriggeredBy: "agent", Action: "STOP", ResourceID: "c2", ResourceType: "cluster", Result: "Failed", Severity: "Error"},
	}

	dtos := ToClusterEventDTOResponseList(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, int64(1), dtos[0].ID)
	assert.Equal(t, int64(2), dtos[1].ID)
}

// TestSystemEventDBResponse_ToSystemEventDTOResponse verifies DB model to DTO conversion.
func TestSystemEventDBResponse_ToSystemEventDTOResponse(t *testing.T) {
	t.Run("Convert SystemEventDBResponse to SystemEventDTOResponse", func(t *testing.T) {
		testSystemEventDBResponse_ToSystemEventDTOResponse_Correct(t)
	})
}

func testSystemEventDBResponse_ToSystemEventDTOResponse_Correct(t *testing.T) {
	now := time.Now().UTC()
	desc := "desc"

	model := SystemEventDBResponse{
		ClusterEventDBResponse: ClusterEventDBResponse{
			ID:             10,
			EventTimestamp: now,
			TriggeredBy:    "scheduler",
			Action:         "START",
			ResourceID:     "cluster-10",
			ResourceType:   "cluster",
			Result:         "Pending",
			Description:    &desc,
			Severity:       "Warning",
		},
		AccountID: "acc-1",
		Provider:  "AWS",
	}

	dto := model.ToSystemEventDTOResponse()

	assert.NotNil(t, dto)
	assert.Equal(t, int64(10), dto.ID)
	assert.Equal(t, now, dto.EventTimestamp)
	assert.Equal(t, "scheduler", dto.TriggeredBy)
	assert.Equal(t, "START", dto.Action)
	assert.Equal(t, "cluster-10", dto.ResourceID)
	assert.Equal(t, "cluster", dto.ResourceType)
	assert.Equal(t, "Pending", dto.Result)
	assert.Equal(t, &desc, dto.Description)
	assert.Equal(t, "Warning", dto.Severity)
	assert.Equal(t, "acc-1", dto.AccountID)
	assert.Equal(t, "AWS", dto.Provider)
}

// TestToSystemEventDTOResponseList verifies list conversion from DB responses to DTOs.
func TestToSystemEventDTOResponseList(t *testing.T) {
	t.Run("Convert SystemEventDBResponse list to DTO list", func(t *testing.T) {
		testToSystemEventDTOResponseList_Correct(t)
	})
}

func testToSystemEventDTOResponseList_Correct(t *testing.T) {
	now := time.Now().UTC()

	models := []SystemEventDBResponse{
		{
			ClusterEventDBResponse: ClusterEventDBResponse{
				ID:             1,
				EventTimestamp: now,
				TriggeredBy:    "api",
				Action:         "START",
				ResourceID:     "c1",
				ResourceType:   "cluster",
				Result:         "Success",
				Severity:       "Info",
			},
			AccountID: "acc-1",
			Provider:  "AWS",
		},
		{
			ClusterEventDBResponse: ClusterEventDBResponse{
				ID:             2,
				EventTimestamp: now.Add(-time.Minute),
				TriggeredBy:    "agent",
				Action:         "STOP",
				ResourceID:     "c2",
				ResourceType:   "cluster",
				Result:         "Failed",
				Severity:       "Error",
			},
			AccountID: "acc-2",
			Provider:  "GCP",
		},
	}

	dtos := ToSystemEventDTOResponseList(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, int64(1), dtos[0].ID)
	assert.Equal(t, "acc-1", dtos[0].AccountID)
	assert.Equal(t, "AWS", dtos[0].Provider)
	assert.Equal(t, int64(2), dtos[1].ID)
	assert.Equal(t, "acc-2", dtos[1].AccountID)
	assert.Equal(t, "GCP", dtos[1].Provider)
}
