package dto

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestInstanceDTORequest_ToInventoryInstance verifies DTO to inventory.Instance conversion.
func TestInstanceDTORequest_ToInventoryInstance(t *testing.T) {
	t.Run("Valid DTO", func(t *testing.T) { testInstanceDTORequest_ToInventoryInstance_Correct(t) })
	t.Run("Invalid DTO returns nil", func(t *testing.T) { testInstanceDTORequest_ToInventoryInstance_Invalid(t) })
}

func testInstanceDTORequest_ToInventoryInstance_Correct(t *testing.T) {
	now := time.Now().UTC()
	createdAt := now.Add(-2 * time.Hour)
	lastScan := now.Add(-30 * time.Minute)

	dto := InstanceDTORequest{
		InstanceID:       "i-123",
		InstanceName:     "node-1",
		InstanceType:     "t3.large",
		Provider:         inventory.AWSProvider,
		AvailabilityZone: "eu-west-1a",
		Status:           inventory.Running,
		ClusterID:        "cluster-1",
		LastScanTS:       lastScan,
		CreatedAt:        createdAt,
		Age:              999,             // not applied by ToInventoryInstance (computed by inventory.NewInstance)
		Owner:            "ignored-owner", // ToInventoryInstance doesn't use Owner
		Tags: []TagDTORequest{
			{Key: "Owner", Value: "team-a"},
			{Key: "env", Value: "prod"},
		},
	}

	instance := dto.ToInventoryInstance()

	assert.NotNil(t, instance)

	assert.Equal(t, dto.InstanceID, instance.InstanceID)
	assert.Equal(t, dto.InstanceName, instance.InstanceName)
	assert.Equal(t, dto.Provider, instance.Provider)
	assert.Equal(t, dto.InstanceType, instance.InstanceType)
	assert.Equal(t, dto.AvailabilityZone, instance.AvailabilityZone)
	assert.Equal(t, dto.Status, instance.Status)
	assert.Equal(t, dto.ClusterID, instance.ClusterID)

	// Tags mapped via ToInventoryTagList / inventory.NewTag, InstanceID should be empty on tags.
	assert.Len(t, instance.Tags, 2)
	assert.Equal(t, "Owner", instance.Tags[0].Key)
	assert.Equal(t, "team-a", instance.Tags[0].Value)
	assert.Equal(t, "", instance.Tags[0].InstanceID)

	assert.Equal(t, lastScan, instance.LastScanTS)
	assert.Equal(t, createdAt, instance.CreatedAt)
}

func testInstanceDTORequest_ToInventoryInstance_Invalid(t *testing.T) {
	dto := InstanceDTORequest{
		InstanceID:       "", // NewInstance should fail
		InstanceName:     "node-1",
		InstanceType:     "t3.large",
		Provider:         inventory.AWSProvider,
		AvailabilityZone: "eu-west-1a",
		Status:           inventory.Running,
		CreatedAt:        time.Now().UTC(),
		Tags:             []TagDTORequest{},
	}

	instance := dto.ToInventoryInstance()
	assert.Nil(t, instance)
}

// TestToInventoryInstanceList verifies slice conversion from DTOs to inventory.Instance.
func TestToInventoryInstanceList(t *testing.T) {
	t.Run("Multiple DTOs", func(t *testing.T) { testToInventoryInstanceList_Correct(t) })
}

func testToInventoryInstanceList_Correct(t *testing.T) {
	now := time.Now().UTC()

	dtos := []InstanceDTORequest{
		{
			InstanceID:       "i-1",
			InstanceName:     "n1",
			InstanceType:     "t3.small",
			Provider:         inventory.AWSProvider,
			AvailabilityZone: "eu-west-1a",
			Status:           inventory.Running,
			ClusterID:        "c1",
			LastScanTS:       now,
			CreatedAt:        now.Add(-time.Hour),
			Tags:             []TagDTORequest{{Key: "env", Value: "dev"}},
		},
		{
			InstanceID:       "i-2",
			InstanceName:     "n2",
			InstanceType:     "t3.medium",
			Provider:         inventory.AWSProvider,
			AvailabilityZone: "eu-west-1b",
			Status:           inventory.Stopped,
			ClusterID:        "c2",
			LastScanTS:       now.Add(-time.Minute),
			CreatedAt:        now.Add(-2 * time.Hour),
			Tags:             []TagDTORequest{{Key: "env", Value: "prod"}},
		},
	}

	instances := ToInventoryInstanceList(dtos)

	assert.NotNil(t, instances)
	assert.Len(t, *instances, 2)
	assert.Equal(t, "i-1", (*instances)[0].InstanceID)
	assert.Equal(t, "i-2", (*instances)[1].InstanceID)
}

// TestToInstanceDTORequest verifies inventory.Instance to DTO conversion.
func TestToInstanceDTORequest(t *testing.T) {
	t.Run("Instance to DTO", func(t *testing.T) { testToInstanceDTORequest_Correct(t) })
}

func testToInstanceDTORequest_Correct(t *testing.T) {
	now := time.Now().UTC()

	instance := inventory.Instance{
		InstanceID:       "i-123",
		InstanceName:     "node-1",
		InstanceType:     "t3.large",
		Provider:         inventory.AWSProvider,
		AvailabilityZone: "eu-west-1a",
		Status:           inventory.Running,
		ClusterID:        "cluster-1",
		LastScanTS:       now,
		CreatedAt:        now.Add(-time.Hour),
		Age:              10,
		Tags: []inventory.Tag{
			{Key: "Owner", Value: "team-a"},
			{Key: "env", Value: "prod"},
		},
	}

	dto := ToInstanceDTORequest(instance)

	assert.NotNil(t, dto)
	assert.Equal(t, instance.InstanceID, dto.InstanceID)
	assert.Equal(t, instance.InstanceName, dto.InstanceName)
	assert.Equal(t, instance.InstanceType, dto.InstanceType)
	assert.Equal(t, instance.Provider, dto.Provider)
	assert.Equal(t, instance.AvailabilityZone, dto.AvailabilityZone)
	assert.Equal(t, instance.Status, dto.Status)
	assert.Equal(t, instance.ClusterID, dto.ClusterID)
	assert.Equal(t, instance.LastScanTS, dto.LastScanTS)
	assert.Equal(t, instance.CreatedAt, dto.CreatedAt)
	assert.Equal(t, instance.Age, dto.Age)

	// Owner is derived from tags
	assert.Equal(t, "team-a", dto.Owner)

	// Tags are mapped back to TagDTORequest list
	assert.Len(t, dto.Tags, 2)
	assert.Equal(t, "Owner", dto.Tags[0].Key)
	assert.Equal(t, "env", dto.Tags[1].Key)
}

// TestToInstanceDTORequestList verifies list conversion from inventory instances to DTO requests.
func TestToInstanceDTORequestList(t *testing.T) {
	t.Run("Instance DTO list", func(t *testing.T) { testToInstanceDTORequestList_Correct(t) })
}

func testToInstanceDTORequestList_Correct(t *testing.T) {
	instances := []inventory.Instance{
		{InstanceID: "i-1", InstanceName: "n1", Provider: inventory.AWSProvider},
		{InstanceID: "i-2", InstanceName: "n2", Provider: inventory.AWSProvider},
	}

	dtoList := ToInstanceDTORequestList(instances)

	assert.NotNil(t, dtoList)
	assert.Len(t, *dtoList, 2)
	assert.Equal(t, "i-1", (*dtoList)[0].InstanceID)
	assert.Equal(t, "i-2", (*dtoList)[1].InstanceID)
}

// TestInstanceDTOResponse_ToInventoryInstance verifies response DTO to inventory.Instance conversion.
func TestInstanceDTOResponse_ToInventoryInstance(t *testing.T) {
	t.Run("Response DTO to inventory.Instance", func(t *testing.T) { testInstanceDTOResponse_ToInventoryInstance_Correct(t) })
}

func testInstanceDTOResponse_ToInventoryInstance_Correct(t *testing.T) {
	now := time.Now().UTC()

	dto := &InstanceDTOResponse{
		InstanceID:       "i-123",
		InstanceName:     "node-1",
		InstanceType:     "t3.large",
		Provider:         inventory.AWSProvider,
		AvailabilityZone: "eu-west-1a",
		Status:           inventory.Running,
		ClusterID:        "cluster-1",
		LastScanTS:       now,
		CreatedAt:        now.Add(-time.Hour),
		Age:              5,
	}

	instance := dto.ToInventoryInstance()

	assert.NotNil(t, instance)
	assert.Equal(t, dto.InstanceID, instance.InstanceID)
	assert.Equal(t, dto.InstanceName, instance.InstanceName)
	assert.Equal(t, dto.InstanceType, instance.InstanceType)
	assert.Equal(t, dto.Provider, instance.Provider)
	assert.Equal(t, dto.AvailabilityZone, instance.AvailabilityZone)
	assert.Equal(t, dto.Status, instance.Status)
	assert.Equal(t, dto.ClusterID, instance.ClusterID)
	assert.Equal(t, dto.LastScanTS, instance.LastScanTS)
	assert.Equal(t, dto.CreatedAt, instance.CreatedAt)
	assert.Equal(t, dto.Age, instance.Age)
}
