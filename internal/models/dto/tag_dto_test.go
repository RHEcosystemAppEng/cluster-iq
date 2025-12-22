package dto

import (
	"testing"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestTagDTORequest_ToInventoryTag verifies DTO to inventory.Tag conversion.
func TestTagDTORequest_ToInventoryTag(t *testing.T) {
	t.Run("Valid DTO", func(t *testing.T) { testTagDTORequest_ToInventoryTag_Correct(t) })
}

func testTagDTORequest_ToInventoryTag_Correct(t *testing.T) {
	dto := TagDTORequest{
		Key:   "env",
		Value: "prod",
	}

	tag := dto.ToInventoryTag()

	assert.NotNil(t, tag)
	assert.Equal(t, dto.Key, tag.Key)
	assert.Equal(t, dto.Value, tag.Value)
	assert.Equal(t, "", tag.InstanceID)
}

// TestToInventoryTagList verifies slice conversion from DTOs to inventory.Tag.
func TestToInventoryTagList(t *testing.T) {
	t.Run("Multiple DTOs", func(t *testing.T) { testToInventoryTagList_Correct(t) })
	t.Run("Empty input", func(t *testing.T) { testToInventoryTagList_Empty(t) })
}

func testToInventoryTagList_Correct(t *testing.T) {
	dtos := []TagDTORequest{
		{Key: "k1", Value: "v1"},
		{Key: "k2", Value: "v2"},
	}

	tags := ToInventoryTagList(dtos)

	assert.NotNil(t, tags)
	assert.Len(t, *tags, 2)
	assert.Equal(t, "k1", (*tags)[0].Key)
	assert.Equal(t, "k2", (*tags)[1].Key)
}

func testToInventoryTagList_Empty(t *testing.T) {
	tags := ToInventoryTagList(nil)

	assert.NotNil(t, tags)
	assert.Len(t, *tags, 0)
}

// TestToTagDTORequest verifies inventory.Tag to DTO conversion.
func TestToTagDTORequest(t *testing.T) {
	t.Run("Tag to DTO", func(t *testing.T) { testToTagDTORequest_Correct(t) })
}

func testToTagDTORequest_Correct(t *testing.T) {
	tag := inventory.Tag{
		Key:   "owner",
		Value: "team-a",
	}

	dto := ToTagDTORequest(tag)

	assert.NotNil(t, dto)
	assert.Equal(t, tag.Key, dto.Key)
	assert.Equal(t, tag.Value, dto.Value)
}

// TestToTagDTORequestList verifies list conversion from inventory tags to DTO requests.
func TestToTagDTORequestList(t *testing.T) {
	t.Run("Tag DTO list", func(t *testing.T) { testToTagDTORequestList_Correct(t) })
	t.Run("Empty input", func(t *testing.T) { testToTagDTORequestList_Empty(t) })
}

func testToTagDTORequestList_Correct(t *testing.T) {
	tags := []inventory.Tag{
		{Key: "k1", Value: "v1"},
		{Key: "k2", Value: "v2"},
	}

	dtoList := ToTagDTORequestList(tags)

	assert.NotNil(t, dtoList)
	assert.Len(t, *dtoList, 2)
	assert.Equal(t, "k1", (*dtoList)[0].Key)
	assert.Equal(t, "k2", (*dtoList)[1].Key)
}

func testToTagDTORequestList_Empty(t *testing.T) {
	dtoList := ToTagDTORequestList(nil)

	assert.NotNil(t, dtoList)
	assert.Len(t, *dtoList, 0)
}
