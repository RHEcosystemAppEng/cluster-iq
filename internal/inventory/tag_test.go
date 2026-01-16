package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewTag verifies NewTag returns a correctly initialized Tag.
func TestNewTag(t *testing.T) {
	t.Run("New Tag", func(t *testing.T) { testNewTag_Correct(t) })
}

func testNewTag_Correct(t *testing.T) {
	key := "env"
	value := "prod"
	instanceID := "i-123"

	tag := NewTag(key, value, instanceID)

	assert.NotNil(t, tag)
	assert.Equal(t, key, tag.Key)
	assert.Equal(t, value, tag.Value)
	assert.Equal(t, instanceID, tag.InstanceID)
}

// TestLookForTagByKey verifies LookForTagByKey returns the expected tag pointer.
func TestLookForTagByKey(t *testing.T) {
	t.Run("Tag found", func(t *testing.T) { testLookForTagByKey_Found(t) })
	t.Run("Tag not found", func(t *testing.T) { testLookForTagByKey_NotFound(t) })
}

func testLookForTagByKey_Found(t *testing.T) {
	tags := []Tag{
		{Key: "Name", Value: "node-1"},
		{Key: "Owner", Value: "alice"},
	}

	res := LookForTagByKey("Owner", tags)
	assert.NotNil(t, res)
	assert.Equal(t, "Owner", res.Key)
	assert.Equal(t, "alice", res.Value)
}

func testLookForTagByKey_NotFound(t *testing.T) {
	tags := []Tag{
		{Key: "Name", Value: "node-1"},
	}

	res := LookForTagByKey("Owner", tags)
	assert.Nil(t, res)
}

// TestGetOwnerFromTags verifies GetOwnerFromTags returns expected values.
func TestGetOwnerFromTags(t *testing.T) {
	t.Run("Owner exists", func(t *testing.T) { testGetOwnerFromTags_Exists(t) })
	t.Run("Owner missing", func(t *testing.T) { testGetOwnerFromTags_Missing(t) })
}

func testGetOwnerFromTags_Exists(t *testing.T) {
	tags := []Tag{
		{Key: "Owner", Value: "john"},
	}

	assert.Equal(t, "john", GetOwnerFromTags(tags))
}

func testGetOwnerFromTags_Missing(t *testing.T) {
	tags := []Tag{
		{Key: "Name", Value: "node-1"},
	}

	assert.Equal(t, "", GetOwnerFromTags(tags))
}

// TestGetInstanceNameFromTags verifies GetInstanceNameFromTags returns expected values.
func TestGetInstanceNameFromTags(t *testing.T) {
	t.Run("Name exists", func(t *testing.T) { testGetInstanceNameFromTags_Exists(t) })
	t.Run("Name missing", func(t *testing.T) { testGetInstanceNameFromTags_Missing(t) })
}

func testGetInstanceNameFromTags_Exists(t *testing.T) {
	tags := []Tag{
		{Key: "Name", Value: "my-instance"},
	}

	assert.Equal(t, "my-instance", GetInstanceNameFromTags(tags))
}

func testGetInstanceNameFromTags_Missing(t *testing.T) {
	tags := []Tag{
		{Key: "Owner", Value: "john"},
	}

	assert.Equal(t, "", GetInstanceNameFromTags(tags))
}

// TestGetClusterNameFromTags verifies GetClusterNameFromTags parses a tag key correctly.
func TestGetClusterNameFromTags(t *testing.T) {
	t.Run("ClusterName found", func(t *testing.T) { testGetClusterNameFromTags_Found(t) })
	t.Run("ClusterName missing", func(t *testing.T) { testGetClusterNameFromTags_Missing(t) })
	t.Run("ClusterName malformed key", func(t *testing.T) { testGetClusterNameFromTags_Malformed(t) })
}

func testGetClusterNameFromTags_Found(t *testing.T) {
	tags := []Tag{
		{Key: ClusterTagKey + "mycluster-ABCDE", Value: "owned"},
	}

	assert.Equal(t, "mycluster", GetClusterNameFromTags(tags))
}

func testGetClusterNameFromTags_Missing(t *testing.T) {
	tags := []Tag{
		{Key: "Name", Value: "node-1"},
	}

	assert.Equal(t, UnknownClusterNameCode, GetClusterNameFromTags(tags))
}

func testGetClusterNameFromTags_Malformed(t *testing.T) {
	tags := []Tag{
		{Key: ClusterTagKey + "invalidkey", Value: "x"},
	}

	assert.Equal(t, UnknownClusterNameCode, GetClusterNameFromTags(tags))
}

// TestGetClusterIDFromTags verifies GetClusterIDFromTags parses a tag key correctly.
func TestGetClusterIDFromTags(t *testing.T) {
	t.Run("ClusterID found", func(t *testing.T) { testGetClusterIDFromTags_Found(t) })
	t.Run("ClusterID missing", func(t *testing.T) { testGetClusterIDFromTags_Missing(t) })
	t.Run("ClusterID malformed key", func(t *testing.T) { testGetClusterIDFromTags_Malformed(t) })
}

func testGetClusterIDFromTags_Found(t *testing.T) {
	tags := []Tag{
		{Key: ClusterTagKey + "mycluster-ABCDE", Value: "owned"},
	}

	assert.Equal(t, "mycluster-ABCDE", GetClusterIDFromTags(tags))
}

func testGetClusterIDFromTags_Missing(t *testing.T) {
	tags := []Tag{
		{Key: "Owner", Value: "john"},
	}

	assert.Equal(t, UnknownClusterNameCode, GetClusterIDFromTags(tags))
}

func testGetClusterIDFromTags_Malformed(t *testing.T) {
	tags := []Tag{
		{Key: ClusterTagKey + "", Value: "x"},
	}

	// parseClusterID returns UnknownClusterIDCode for non-matching keys.
	assert.Equal(t, UnknownClusterIDCode, GetClusterIDFromTags(tags))
}

// TestGetInfraIDFromTags verifies GetInfraIDFromTags parses infraID correctly.
func TestGetInfraIDFromTags(t *testing.T) {
	t.Run("InfraID found", func(t *testing.T) { testGetInfraIDFromTags_Found(t) })
	t.Run("InfraID missing", func(t *testing.T) { testGetInfraIDFromTags_Missing(t) })
	t.Run("InfraID malformed key", func(t *testing.T) { testGetInfraIDFromTags_Malformed(t) })
}

func testGetInfraIDFromTags_Found(t *testing.T) {
	tags := []Tag{
		{Key: ClusterTagKey + "/mycluster-ABCDE", Value: "owned"},
	}

	assert.Equal(t, "ABCDE", GetInfraIDFromTags(tags))
}

func testGetInfraIDFromTags_Missing(t *testing.T) {
	tags := []Tag{
		{Key: "Owner", Value: "john"},
	}

	assert.Equal(t, UnknownClusterNameCode, GetInfraIDFromTags(tags))
}

func testGetInfraIDFromTags_Malformed(t *testing.T) {
	tags := []Tag{
		{Key: ClusterTagKey + "/invalidkey", Value: "x"},
	}

	// parseInfraID returns empty string for non-matching keys, while GetInfraIDFromTags
	// returns that value once it finds ClusterTagKey.
	assert.Equal(t, "", GetInfraIDFromTags(tags))
}
