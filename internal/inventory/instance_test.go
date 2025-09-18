package inventory

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: Group by function and test
// TODO: Include Asserts

// TestNewInstance verifies that NewInstance returns a correctly initialized instance
func TestNewInstance(t *testing.T) {
	id := InstanceID("0000-11A")
	name := "testAccount"
	var provider Provider = UnknownProvider
	instanceType := "t2.micro"
	availabilityZone := "us-west-1a"
	status := Terminated
	clusterID := "testCluster"
	tags := make([]Tag, 0)
	createdAt := time.Now()

	expectedInstance := &Instance{
		InstanceID:       id,
		InstanceName:     name,
		Provider:         provider,
		InstanceType:     instanceType,
		AvailabilityZone: availabilityZone,
		Status:           status,
		ClusterID:        clusterID,
		LastScanTS:       createdAt,
		CreatedAt:        createdAt,
		Tags:             tags,
	}

	actualInstance := NewInstance(id, name, provider, instanceType, availabilityZone, status, tags, createdAt)

	assert.NotNil(t, actualInstance)
	assert.NotZero(t, actualInstance.LastScanTS)

	expectedInstance.Age = actualInstance.Age
	expectedInstance.LastScanTS = actualInstance.LastScanTS
	expectedInstance.CreatedAt = actualInstance.CreatedAt
	assert.Equal(t, expectedInstance, actualInstance)
}

// TestAddTag verifies that a new tag is appended to the tag list
func TestAddTag(t *testing.T) {
	i := Instance{}
	tag := Tag{Key: "env", Value: "prod"}
	i.AddTag(tag)
	if len(i.Tags) != 1 || i.Tags[0] != tag {
		t.Errorf("tag was not added correctly")
	}
}

// TestInstance_String verifies String method returns expected format
func TestInstance_String(t *testing.T) {
	i := Instance{
		InstanceID:       "i-123",
		InstanceName:     "test",
		Provider:         AWSProvider,
		InstanceType:     "t2.micro",
		AvailabilityZone: "us-east-1a",
		Status:           Running,
		ClusterID:        "cluster-x",
		Expenses:         []Expense{{Amount: 5}},
	}

	str := i.String()
	if !(strings.Contains(str, "test") && strings.Contains(str, "AWS") && strings.Contains(str, "t2.micro")) {
		t.Errorf("unexpected output from String(): %s", str)
	}
}

// TestPrintInstance verifies PrintInstance runs without panic
func TestPrintInstance(t *testing.T) {
	i := Instance{InstanceID: "i-456", InstanceName: "node1"}
	i.PrintInstance()
}
