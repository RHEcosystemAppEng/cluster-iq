package db_test

import (
	"encoding/json"
	"testing"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/convert"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/stretchr/testify/assert"
)

// TestTagDBResponse_ToTagDTOResponse verifies DB model to DTO conversion.
func TestTagDBResponse_ToTagDTOResponse(t *testing.T) {
	t.Run("Convert TagDBResponse to TagDTOResponse", func(t *testing.T) {
		testTagDBResponse_ToTagDTOResponse_Correct(t)
	})
}

func testTagDBResponse_ToTagDTOResponse_Correct(t *testing.T) {
	conv := &convert.ConverterImpl{}

	model := db.TagDBResponse{
		Key:        "env",
		Value:      "prod",
		InstanceID: "i-123",
	}

	dto := conv.ToTagDTO(model)

	assert.Equal(t, model.Key, dto.Key)
	assert.Equal(t, model.Value, dto.Value)
	assert.Equal(t, model.InstanceID, dto.InstanceID)
}

// TestToTagsDTOResponseList verifies slice conversion from DB responses to DTOs.
func TestToTagsDTOResponseList(t *testing.T) {
	t.Run("Convert TagDBResponse list to DTO list", func(t *testing.T) {
		testToTagsDTOResponseList_Correct(t)
	})
}

func testToTagsDTOResponseList_Correct(t *testing.T) {
	conv := &convert.ConverterImpl{}

	models := []db.TagDBResponse{
		{Key: "k1", Value: "v1", InstanceID: "i-1"},
		{Key: "k2", Value: "v2", InstanceID: "i-2"},
	}

	dtos := conv.ToTagDTOs(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, "k1", dtos[0].Key)
	assert.Equal(t, "k2", dtos[1].Key)
}

// TestTagDBResponses_Scan verifies Scan behavior for all supported inputs.
func TestTagDBResponses_Scan(t *testing.T) {
	t.Run("Scan nil value", func(t *testing.T) { testTagDBResponsesScan_Nil(t) })
	t.Run("Scan JSON array", func(t *testing.T) { testTagDBResponsesScan_Array(t) })
	t.Run("Scan JSON object", func(t *testing.T) { testTagDBResponsesScan_Map(t) })
	t.Run("Scan invalid JSON", func(t *testing.T) { testTagDBResponsesScan_InvalidJSON(t) })
}

func testTagDBResponsesScan_Nil(t *testing.T) {
	var tags db.TagDBResponses

	err := tags.Scan(nil)

	assert.NoError(t, err)
	assert.Len(t, tags, 0)
}

func testTagDBResponsesScan_Array(t *testing.T) {
	input := []db.TagDBResponse{
		{Key: "env", Value: "prod"},
		{Key: "team", Value: "platform"},
	}
	b, _ := json.Marshal(input)

	var tags db.TagDBResponses
	err := tags.Scan(b)

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "env", tags[0].Key)
	assert.Equal(t, "team", tags[1].Key)
}

func testTagDBResponsesScan_Map(t *testing.T) {
	input := map[string]interface{}{
		"env":  "prod",
		"cost": 10,
	}
	b, _ := json.Marshal(input)

	var tags db.TagDBResponses
	err := tags.Scan(string(b))

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
}

func testTagDBResponsesScan_InvalidJSON(t *testing.T) {
	var tags db.TagDBResponses

	err := tags.Scan([]byte("{invalid-json"))

	assert.Error(t, err)
}
