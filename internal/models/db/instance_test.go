package db_test

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/convert"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/stretchr/testify/assert"
)

// TestInstanceDBResponse_ToInstanceDTOResponse verifies DB model to DTO conversion.
func TestInstanceDBResponse_ToInstanceDTOResponse(t *testing.T) {
	t.Run("Convert InstanceDBResponse to InstanceDTOResponse", func(t *testing.T) {
		testInstanceDBResponse_ToInstanceDTOResponse_Correct(t)
	})
}

func testInstanceDBResponse_ToInstanceDTOResponse_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	model := db.InstanceDBResponse{
		InstanceID:            "i-123",
		InstanceName:          "node-1",
		InstanceType:          "t3.large",
		Provider:              inventory.AWSProvider,
		AvailabilityZone:      "eu-west-1a",
		Status:                inventory.Running,
		ClusterID:             "cluster-1",
		ClusterName:           "cluster-name",
		LastScanTimestamp:     now,
		CreatedAt:             now.Add(-time.Hour),
		Age:                   5,
		TotalCost:             123.45,
		Last15DaysCost:        12.34,
		LastMonthCost:         56.78,
		CurrentMonthSoFarCost: 9.01,
		Tags: db.TagDBResponses{
			{Key: "Owner", Value: "team-a", InstanceID: "i-123"},
			{Key: "env", Value: "prod", InstanceID: "i-123"},
		},
	}

	dto := conv.ToInstanceDTO(model)

	assert.Equal(t, model.InstanceID, dto.InstanceID)
	assert.Equal(t, model.InstanceName, dto.InstanceName)
	assert.Equal(t, model.InstanceType, dto.InstanceType)
	assert.Equal(t, model.Provider, dto.Provider)
	assert.Equal(t, model.AvailabilityZone, dto.AvailabilityZone)
	assert.Equal(t, model.Status, dto.Status)
	assert.Equal(t, model.ClusterID, dto.ClusterID)
	assert.Equal(t, model.ClusterName, dto.ClusterName)
	assert.Equal(t, model.LastScanTimestamp, dto.LastScanTimestamp)
	assert.Equal(t, model.CreatedAt, dto.CreatedAt)
	assert.Equal(t, model.Age, dto.Age)
	assert.Equal(t, model.TotalCost, dto.TotalCost)
	assert.Equal(t, model.Last15DaysCost, dto.Last15DaysCost)
	assert.Equal(t, model.LastMonthCost, dto.LastMonthCost)
	assert.Equal(t, model.CurrentMonthSoFarCost, dto.CurrentMonthSoFarCost)

	assert.Len(t, dto.Tags, 2)
	assert.Equal(t, "Owner", dto.Tags[0].Key)
	assert.Equal(t, "env", dto.Tags[1].Key)
}

// TestToInstanceDTOResponseList verifies list conversion from DB responses to DTOs.
func TestToInstanceDTOResponseList(t *testing.T) {
	t.Run("Convert InstanceDBResponse list to DTO list", func(t *testing.T) {
		testToInstanceDTOResponseList_Correct(t)
	})
}

func testToInstanceDTOResponseList_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	models := []db.InstanceDBResponse{
		{
			InstanceID:        "i-1",
			InstanceName:      "n1",
			Provider:          inventory.AWSProvider,
			LastScanTimestamp: now,
			CreatedAt:         now.Add(-time.Hour),
		},
		{
			InstanceID:        "i-2",
			InstanceName:      "n2",
			Provider:          inventory.AWSProvider,
			LastScanTimestamp: now.Add(-time.Minute),
			CreatedAt:         now.Add(-2 * time.Hour),
		},
	}

	dtos := conv.ToInstanceDTOs(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, "i-1", dtos[0].InstanceID)
	assert.Equal(t, "i-2", dtos[1].InstanceID)
}
