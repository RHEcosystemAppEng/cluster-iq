package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// TestActionDBResponse_ToActionDTOResponse verifies DB model to DTO conversion.
func TestActionDBResponse_ToActionDTOResponse(t *testing.T) {
	t.Run("Convert with valid time and cron", func(t *testing.T) { testActionDBResponse_ToActionDTOResponse_WithValidFields(t) })
	t.Run("Convert with invalid time and cron", func(t *testing.T) { testActionDBResponse_ToActionDTOResponse_WithInvalidFields(t) })
}

func testActionDBResponse_ToActionDTOResponse_WithValidFields(t *testing.T) {
	now := time.Now().UTC()

	model := ActionDBResponse{
		ID:        "id-1",
		Type:      "scheduled_action",
		Time:      sql.NullTime{Time: now, Valid: true},
		CronExp:   sql.NullString{String: "0 0 * * *", Valid: true},
		Operation: "START",
		Status:    "Pending",
		Enabled:   true,
		ClusterID: "cluster-1",
		Region:    "eu-west-1",
		AccountID: "acc-1",
		Instances: pq.StringArray{"i-1", "i-2"},
	}

	dto := model.ToActionDTOResponse()

	assert.NotNil(t, dto)
	assert.Equal(t, model.ID, dto.ID)
	assert.Equal(t, model.Type, dto.Type)
	assert.Equal(t, now, dto.Time)
	assert.Equal(t, "0 0 * * *", dto.CronExp)
	assert.Equal(t, model.Operation, dto.Operation)
	assert.Equal(t, model.Status, dto.Status)
	assert.Equal(t, model.Enabled, dto.Enabled)
	assert.Equal(t, model.ClusterID, dto.ClusterID)
	assert.Equal(t, model.Region, dto.Region)
	assert.Equal(t, model.AccountID, dto.AccountID)
	assert.Equal(t, []string{"i-1", "i-2"}, dto.Instances)
}

func testActionDBResponse_ToActionDTOResponse_WithInvalidFields(t *testing.T) {
	model := ActionDBResponse{
		ID:        "id-2",
		Type:      "cron_action",
		Time:      sql.NullTime{Valid: false},
		CronExp:   sql.NullString{Valid: false},
		Operation: "STOP",
		Status:    "Failed",
		Enabled:   false,
		ClusterID: "cluster-2",
		Region:    "us-east-1",
		AccountID: "acc-2",
		Instances: pq.StringArray{"i-9"},
	}

	dto := model.ToActionDTOResponse()

	assert.NotNil(t, dto)
	assert.Equal(t, model.ID, dto.ID)
	assert.Equal(t, model.Type, dto.Type)

	// When invalid, dto.Time is the zero value
	assert.True(t, dto.Time.IsZero())
	assert.Equal(t, "", dto.CronExp)

	assert.Equal(t, model.Operation, dto.Operation)
	assert.Equal(t, model.Status, dto.Status)
	assert.Equal(t, model.Enabled, dto.Enabled)
	assert.Equal(t, model.ClusterID, dto.ClusterID)
	assert.Equal(t, model.Region, dto.Region)
	assert.Equal(t, model.AccountID, dto.AccountID)
	assert.Equal(t, []string{"i-9"}, dto.Instances)
}

// TestToActionDTOResponseList verifies list conversion from DB responses to DTOs.
func TestToActionDTOResponseList(t *testing.T) {
	t.Run("Convert ActionDBResponse list to DTO list", func(t *testing.T) { testToActionDTOResponseList_Correct(t) })
}

func testToActionDTOResponseList_Correct(t *testing.T) {
	now := time.Now().UTC()

	models := []ActionDBResponse{
		{
			ID:        "id-1",
			Type:      "scheduled_action",
			Time:      sql.NullTime{Time: now, Valid: true},
			CronExp:   sql.NullString{String: "", Valid: false},
			Operation: "START",
			Status:    "Pending",
			Enabled:   true,
			ClusterID: "cluster-1",
			Region:    "eu-west-1",
			AccountID: "acc-1",
			Instances: pq.StringArray{"i-1"},
		},
		{
			ID:        "id-2",
			Type:      "cron_action",
			Time:      sql.NullTime{Valid: false},
			CronExp:   sql.NullString{String: "*/5 * * * *", Valid: true},
			Operation: "STOP",
			Status:    "Running",
			Enabled:   false,
			ClusterID: "cluster-2",
			Region:    "us-east-1",
			AccountID: "acc-2",
			Instances: pq.StringArray{"i-2"},
		},
	}

	dtos := ToActionDTOResponseList(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, "id-1", dtos[0].ID)
	assert.Equal(t, "id-2", dtos[1].ID)
	assert.Equal(t, "scheduled_action", dtos[0].Type)
	assert.Equal(t, "cron_action", dtos[1].Type)
}
