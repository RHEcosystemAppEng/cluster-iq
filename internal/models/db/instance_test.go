package db

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestTagDBResponseJSON_ToTagDTOResponse verifies JSON tag model to DTO conversion.
func TestTagDBResponseJSON_ToTagDTOResponse(t *testing.T) {
	t.Run("Convert TagDBResponseJSON to TagDTOResponse", func(t *testing.T) {
		testTagDBResponseJSON_ToTagDTOResponse_Correct(t)
	})
}

func testTagDBResponseJSON_ToTagDTOResponse_Correct(t *testing.T) {
	tag := TagDBResponseJSON{Key: "env", Value: "prod"}

	dto := tag.ToTagDTOResponse()

	assert.NotNil(t, dto)
	assert.Equal(t, "env", dto.Key)
	assert.Equal(t, "prod", dto.Value)
}

// TestTagDBResponseList_Scan verifies Scan behavior for all supported inputs.
func TestTagDBResponseList_Scan(t *testing.T) {
	t.Run("Scan nil value", func(t *testing.T) { testTagDBResponseList_Scan_Nil(t) })
	t.Run("Scan non []byte value", func(t *testing.T) { testTagDBResponseList_Scan_NonBytes(t) })
	t.Run("Scan valid JSON", func(t *testing.T) { testTagDBResponseList_Scan_ValidJSON(t) })
	t.Run("Scan invalid JSON", func(t *testing.T) { testTagDBResponseList_Scan_InvalidJSON(t) })
}

func testTagDBResponseList_Scan_Nil(t *testing.T) {
	var tags TagDBResponseList

	err := tags.Scan(nil)

	assert.NoError(t, err)
	assert.Nil(t, tags)
}

func testTagDBResponseList_Scan_NonBytes(t *testing.T) {
	var tags TagDBResponseList

	err := tags.Scan("not-bytes")

	assert.Error(t, err)
	assert.Nil(t, tags)
}

func testTagDBResponseList_Scan_ValidJSON(t *testing.T) {
	input := []TagDBResponseJSON{
		{Key: "k1", Value: "v1"},
		{Key: "k2", Value: "v2"},
	}
	b, _ := json.Marshal(input)

	var tags TagDBResponseList
	err := tags.Scan(b)

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "k1", tags[0].Key)
	assert.Equal(t, "k2", tags[1].Key)
}

func testTagDBResponseList_Scan_InvalidJSON(t *testing.T) {
	var tags TagDBResponseList

	err := tags.Scan([]byte("{invalid-json"))

	assert.Error(t, err)
	assert.Nil(t, tags)
}

// TestTagDBResponseList_ToTagDTOResponseList verifies conversion from tag list to DTO list.
func TestTagDBResponseList_ToTagDTOResponseList(t *testing.T) {
	t.Run("Convert tag list to DTO list", func(t *testing.T) { testTagDBResponseList_ToTagDTOResponseList_Correct(t) })
	t.Run("Convert empty tag list", func(t *testing.T) { testTagDBResponseList_ToTagDTOResponseList_Empty(t) })
}

func testTagDBResponseList_ToTagDTOResponseList_Correct(t *testing.T) {
	tags := TagDBResponseList{
		{Key: "env", Value: "prod"},
		{Key: "team", Value: "platform"},
	}

	dtos := tags.ToTagDTOResponseList()

	assert.NotNil(t, dtos)
	assert.Len(t, *dtos, 2)
	assert.Equal(t, "env", (*dtos)[0].Key)
	assert.Equal(t, "team", (*dtos)[1].Key)
}

func testTagDBResponseList_ToTagDTOResponseList_Empty(t *testing.T) {
	tags := TagDBResponseList{}

	dtos := tags.ToTagDTOResponseList()

	assert.NotNil(t, dtos)
	assert.Len(t, *dtos, 0)
}

// TestInstanceDBResponse_ToInstanceDTOResponse verifies DB model to DTO conversion.
func TestInstanceDBResponse_ToInstanceDTOResponse(t *testing.T) {
	t.Run("Convert InstanceDBResponse to InstanceDTOResponse", func(t *testing.T) {
		testInstanceDBResponse_ToInstanceDTOResponse_Correct(t)
	})
}

func testInstanceDBResponse_ToInstanceDTOResponse_Correct(t *testing.T) {
	now := time.Now().UTC()

	model := InstanceDBResponse{
		InstanceID:            "i-123",
		InstanceName:          "node-1",
		InstanceType:          "t3.large",
		Provider:              inventory.AWSProvider,
		AvailabilityZone:      "eu-west-1a",
		Status:                inventory.Running,
		ClusterID:             "cluster-1",
		ClusterName:           "cluster-name",
		LastScanTS:            now,
		CreatedAt:             now.Add(-time.Hour),
		Age:                   5,
		TotalCost:             123.45,
		Last15DaysCost:        12.34,
		LastMonthCost:         56.78,
		CurrentMonthSoFarCost: 9.01,
		Tags: TagDBResponses{
			{Key: "Owner", Value: "team-a", InstanceID: "i-123"},
			{Key: "env", Value: "prod", InstanceID: "i-123"},
		},
	}

	dto := model.ToInstanceDTOResponse()

	assert.NotNil(t, dto)
	assert.Equal(t, model.InstanceID, dto.InstanceID)
	assert.Equal(t, model.InstanceName, dto.InstanceName)
	assert.Equal(t, model.InstanceType, dto.InstanceType)
	assert.Equal(t, model.Provider, dto.Provider)
	assert.Equal(t, model.AvailabilityZone, dto.AvailabilityZone)
	assert.Equal(t, model.Status, dto.Status)
	assert.Equal(t, model.ClusterID, dto.ClusterID)
	assert.Equal(t, model.ClusterName, dto.ClusterName)
	assert.Equal(t, model.LastScanTS, dto.LastScanTS)
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

	models := []InstanceDBResponse{
		{
			InstanceID:   "i-1",
			InstanceName: "n1",
			Provider:     inventory.AWSProvider,
			LastScanTS:   now,
			CreatedAt:    now.Add(-time.Hour),
		},
		{
			InstanceID:   "i-2",
			InstanceName: "n2",
			Provider:     inventory.AWSProvider,
			LastScanTS:   now.Add(-time.Minute),
			CreatedAt:    now.Add(-2 * time.Hour),
		},
	}

	dtos := ToInstanceDTOResponseList(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, "i-1", dtos[0].InstanceID)
	assert.Equal(t, "i-2", dtos[1].InstanceID)
}
