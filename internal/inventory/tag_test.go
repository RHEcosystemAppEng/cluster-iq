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
		t.Errorf("Expected tag with key Owner, got %v", tag)
	}
}

// TestLookForTagByKey_NotFound verifies nil is returned if key doesn't exist
func TestLookForTagByKey_NotFound(t *testing.T) {
	tags := []Tag{{Key: "Name", Value: "test"}}
	if tag := LookForTagByKey("Zone", tags); tag != nil {
		t.Errorf("Expected nil, got %v", tag)
	}
}

// TestParseClusterName_Valid extracts name from valid tag key
func TestParseClusterName_Valid(t *testing.T) {
	expectedValue := "test-cluster"
	key := "kubernetes.io/cluster/" + expectedValue + "-abcde"
	actualValue := parseClusterName(key)
	if expectedValue != actualValue {
		t.Errorf("Expected %s, got %s", expectedValue, actualValue)
	}
}

// TestParseClusterName_Invalid returns UNKNOWN when tag is malformed
func TestParseClusterName_Invalid(t *testing.T) {
	key := "kubernetes.io/cluster/invalid"
	name := parseClusterName(key)
	if name != UnknownClusterNameCode {
		t.Errorf("Expected UNKNOWN, got %s", name)
	}
}

// TestParseClusterID_Valid tests clusterID regex logic with valid input
func TestParseClusterID_Valid(t *testing.T) {
	expectedValue := "test-cluster-abcde"
	key := "kubernetes.io/cluster/" + expectedValue
	actualValue := parseClusterID(key)
	if actualValue != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, actualValue)
	}
}

// TestParseClusterID_Invalid tests fallback behavior on malformed key
func TestParseClusterID_Invalid(t *testing.T) {
	key := "garbage"
	id := parseClusterID(key)
	if id != UnknownClusterIDCode {
		t.Errorf("Expected UNKNOWN, got %s", id)
	}
}

// TestParseInfraID_Valid verifies infra ID parsing for valid keys
func TestParseInfraID_Valid(t *testing.T) {
	expectedValue := "abcde"
	key := "kubernetes.io/cluster/test-cluster-" + expectedValue
	actualValue := parseInfraID(key)
	if expectedValue != actualValue {
		t.Errorf("Expected %s, got %s", expectedValue, actualValue)
	}
}

// TestParseInfraID_Invalid returns empty string when no match
func TestParseInfraID_Invalid(t *testing.T) {
	key := "garbage"
	infra := parseInfraID(key)
	if infra != "" {
		t.Errorf("Expected empty string, got %s", infra)
	}
}

// TestGetOwnerFromTags returns Owner tag value if present
func TestGetOwnerFromTags(t *testing.T) {
	expectedValue := "Alice"
	tags := []Tag{{Key: "Owner", Value: expectedValue}}
	actualValue := GetOwnerFromTags(tags)

	if actualValue != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, actualValue)
	}
}

// TestGetOwnerFromTags_Empty if no Owner tag present
func TestGetOwnerFromTags_Empty(t *testing.T) {
	tags := []Tag{{Key: "Name", Value: "node"}}
	val := GetOwnerFromTags(tags)

	if val != "" {
		t.Errorf("Expected empty string, got %s", val)
	}
}

// TestGetInstanceNameFromTags returns Name tag if present
func TestGetInstanceNameFromTags(t *testing.T) {
	expectedValue := "instance-001"
	tags := []Tag{{Key: "Name", Value: expectedValue}}
	actualValue := GetInstanceNameFromTags(tags)

	if actualValue != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, actualValue)
	}
}

// TestGetInstanceNameFromTags_Empty validates fallback to empty string
func TestGetInstanceNameFromTags_Empty(t *testing.T) {
	tags := []Tag{{Key: "Owner", Value: "ops"}}
	val := GetInstanceNameFromTags(tags)
	if val != "" {
		t.Errorf("Expected empty string, got %s", val)
	}
}

// TestGetClusterIDFromTags parses valid cluster tag key
func TestGetClusterIDFromTags(t *testing.T) {
	expectedValue := "test-cluster-abcde"
	tags := []Tag{{Key: "kubernetes.io/cluster/" + expectedValue}}
	actualValue := GetClusterIDFromTags(tags)

	if actualValue != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, actualValue)
	}
}

// TestGetClusterIDFromTags_Unknown when tag is not found
func TestGetClusterIDFromTags_Unknown(t *testing.T) {
	id := GetClusterIDFromTags([]Tag{})
	if id != UnknownClusterNameCode {
		t.Errorf("Expected UNKNOWN, got %s", id)
	}
}

// TestGetClusterNameFromTags works with proper tag key
func TestGetClusterNameFromTags(t *testing.T) {
	expectedValue := "test-cluster"
	tags := []Tag{{Key: "kubernetes.io/cluster/" + expectedValue + "-abcde"}}
	actualValue := GetClusterNameFromTags(tags)

	if actualValue != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, actualValue)
	}
}

// TestGetClusterNameFromTags_Unknown fallback logic
func TestGetClusterNameFromTags_Unknown(t *testing.T) {
	name := GetClusterNameFromTags([]Tag{})
	if name != UnknownClusterNameCode {
		t.Errorf("Expected UNKNOWN, got %s", name)
	}
}

// TestGetInfraIDFromTags verifies successful parsing
func TestGetInfraIDFromTags(t *testing.T) {
	tags := []Tag{{Key: "kubernetes.io/cluster/test-name-abcde"}}
	infra := GetInfraIDFromTags(tags)
	if infra != "abcde" {
		t.Errorf("Expected abcde, got %s", infra)
	}
}

// TestGetInfraIDFromTags_Unknown fallback behavior
func TestGetInfraIDFromTags_Unknown(t *testing.T) {
	infra := GetInfraIDFromTags([]Tag{})
	if infra != UnknownClusterNameCode {
		t.Errorf("Expected UNKNOWN, got %s", infra)
	}
}
