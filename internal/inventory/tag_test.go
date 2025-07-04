package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewTag verifies NewTag constructor assigns values correctly
func TestNewTag(t *testing.T) {
	key := "environment"
	value := "production"
	instanceID := "testInstance"

	expectedTag := &Tag{
		Key:        key,
		Value:      value,
		InstanceID: instanceID,
	}

	actualTag := NewTag(key, value, instanceID)

	assert.NotNil(t, actualTag)
	assert.Equal(t, actualTag, expectedTag)
}

// TestLookForTagByKey_Found verifies tag retrieval by key
func TestLookForTagByKey_Found(t *testing.T) {
	tags := []Tag{
		{Key: "Name", Value: "test"},
		{Key: "Owner", Value: "devops"},
	}
	tag := LookForTagByKey("Owner", tags)
	if tag == nil || tag.Key != "Owner" {
		t.Errorf("expected tag with key Owner, got %v", tag)
	}
}

// TestLookForTagByKey_NotFound verifies nil is returned if key doesn't exist
func TestLookForTagByKey_NotFound(t *testing.T) {
	tags := []Tag{{Key: "Name", Value: "test"}}
	if tag := LookForTagByKey("Zone", tags); tag != nil {
		t.Errorf("expected nil, got %v", tag)
	}
}

// TestParseClusterName_Valid extracts name from valid tag key
func TestParseClusterName_Valid(t *testing.T) {
	key := "kubernetes.io/cluster/ci-ln-abcde"
	name := parseClusterName(key)
	if name != "ci-ln" {
		t.Errorf("expected ci-ln, got %s", name)
	}
}

// TestParseClusterName_Invalid returns UNKNOWN when tag is malformed
func TestParseClusterName_Invalid(t *testing.T) {
	key := "kubernetes.io/cluster/invalid"
	name := parseClusterName(key)
	if name != UnknownClusterNameCode {
		t.Errorf("expected UNKNOWN, got %s", name)
	}
}

// TestParseClusterID_Valid tests clusterID regex logic with valid input
func TestParseClusterID_Valid(t *testing.T) {
	key := "kubernetes.io/cluster/test-name-abcde"
	id := parseClusterID(key)
	if id != "test-name" {
		t.Errorf("expected test-name, got %s", id)
	}
}

// TestParseClusterID_Invalid tests fallback behavior on malformed key
func TestParseClusterID_Invalid(t *testing.T) {
	key := "garbage"
	id := parseClusterID(key)
	if id != UnknownClusterIDCode {
		t.Errorf("expected UNKNOWN, got %s", id)
	}
}

// TestParseInfraID_Valid verifies infra ID parsing for valid keys
func TestParseInfraID_Valid(t *testing.T) {
	key := "kubernetes.io/cluster/test-abcde"
	infra := parseInfraID(key)
	if infra != "abcde" {
		t.Errorf("expected abcde, got %s", infra)
	}
}

// TestParseInfraID_Invalid returns empty string when no match
func TestParseInfraID_Invalid(t *testing.T) {
	key := "something-else"
	infra := parseInfraID(key)
	if infra != "" {
		t.Errorf("expected empty string, got %s", infra)
	}
}

// TestGetOwnerFromTags returns Owner tag value if present
func TestGetOwnerFromTags(t *testing.T) {
	tags := []Tag{{Key: "Owner", Value: "alice"}}
	val := GetOwnerFromTags(tags)
	if val != "Owner" {
		t.Errorf("expected Owner, got %s", val)
	}
}

// TestGetOwnerFromTags_Empty if no Owner tag present
func TestGetOwnerFromTags_Empty(t *testing.T) {
	tags := []Tag{{Key: "Name", Value: "node"}}
	val := GetOwnerFromTags(tags)
	if val != "" {
		t.Errorf("expected empty string, got %s", val)
	}
}

// TestGetInstanceNameFromTags returns Name tag if present
func TestGetInstanceNameFromTags(t *testing.T) {
	tags := []Tag{{Key: "Name", Value: "node-1"}}
	val := GetInstanceNameFromTags(tags)
	if val != "Name" {
		t.Errorf("expected Name, got %s", val)
	}
}

// TestGetInstanceNameFromTags_Empty validates fallback to empty string
func TestGetInstanceNameFromTags_Empty(t *testing.T) {
	tags := []Tag{{Key: "Owner", Value: "ops"}}
	val := GetInstanceNameFromTags(tags)
	if val != "" {
		t.Errorf("expected empty string, got %s", val)
	}
}

// TestGetClusterIDFromTags parses valid cluster tag key
func TestGetClusterIDFromTags(t *testing.T) {
	tags := []Tag{{Key: "kubernetes.io/cluster/test-name-abcde"}}
	id := GetClusterIDFromTags(tags)
	if id != "test-name" {
		t.Errorf("expected test-name, got %s", id)
	}
}

// TestGetClusterIDFromTags_Unknown when tag is not found
func TestGetClusterIDFromTags_Unknown(t *testing.T) {
	id := GetClusterIDFromTags([]Tag{})
	if id != UnknownClusterNameCode {
		t.Errorf("expected UNKNOWN, got %s", id)
	}
}

// TestGetClusterNameFromTags works with proper tag key
func TestGetClusterNameFromTags(t *testing.T) {
	tags := []Tag{{Key: "kubernetes.io/cluster/test-name-abcde"}}
	name := GetClusterNameFromTags(tags)
	if name != "test-name" {
		t.Errorf("expected test-name, got %s", name)
	}
}

// TestGetClusterNameFromTags_Unknown fallback logic
func TestGetClusterNameFromTags_Unknown(t *testing.T) {
	name := GetClusterNameFromTags([]Tag{})
	if name != UnknownClusterNameCode {
		t.Errorf("expected UNKNOWN, got %s", name)
	}
}

// TestGetInfraIDFromTags verifies successful parsing
func TestGetInfraIDFromTags(t *testing.T) {
	tags := []Tag{{Key: "kubernetes.io/cluster/test-name-abcde"}}
	infra := GetInfraIDFromTags(tags)
	if infra != "abcde" {
		t.Errorf("expected abcde, got %s", infra)
	}
}

// TestGetInfraIDFromTags_Unknown fallback behavior
func TestGetInfraIDFromTags_Unknown(t *testing.T) {
	infra := GetInfraIDFromTags([]Tag{})
	if infra != UnknownClusterNameCode {
		t.Errorf("expected UNKNOWN, got %s", infra)
	}
}
