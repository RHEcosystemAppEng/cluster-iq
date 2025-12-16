package dto

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestClusterDTORequest_ToInventoryCluster verifies DTO to inventory.Cluster conversion.
func TestClusterDTORequest_ToInventoryCluster(t *testing.T) {
	t.Run("Valid DTO", func(t *testing.T) { testClusterDTORequest_ToInventoryCluster_Correct(t) })
	t.Run("Invalid DTO returns nil", func(t *testing.T) { testClusterDTORequest_ToInventoryCluster_Invalid(t) })
}

func testClusterDTORequest_ToInventoryCluster_Correct(t *testing.T) {
	now := time.Now().UTC()

	dto := ClusterDTORequest{
		ClusterID:   "cluster-ignored", // NewCluster generates ClusterID; DTO's ClusterID is not used in ToInventoryCluster
		ClusterName: "testCluster",
		InfraID:     "ABCDE",
		Provider:    inventory.AWSProvider,
		Status:      inventory.Running,
		Region:      "eu-west-1",
		AccountID:   "acc-1",
		ConsoleLink: "https://console",
		LastScanTS:  now,
		CreatedAt:   now.Add(-2 * time.Hour),
		Age:         123, // not applied by ToInventoryCluster
		Owner:       "owner",
	}

	cluster := dto.ToInventoryCluster()

	assert.NotNil(t, cluster)

	assert.Equal(t, dto.ClusterName, cluster.ClusterName)
	assert.Equal(t, dto.InfraID, cluster.InfraID)
	assert.Equal(t, dto.Provider, cluster.Provider)
	assert.Equal(t, dto.Region, cluster.Region)
	assert.Equal(t, dto.ConsoleLink, cluster.ConsoleLink)
	assert.Equal(t, dto.Owner, cluster.Owner)

	assert.Equal(t, dto.LastScanTS, cluster.LastScanTS)
	assert.Equal(t, dto.CreatedAt, cluster.CreatedAt)
	assert.Equal(t, dto.Status, cluster.Status)
	assert.Equal(t, dto.AccountID, cluster.AccountID)

	assert.NotEmpty(t, cluster.ClusterID)
}

func testClusterDTORequest_ToInventoryCluster_Invalid(t *testing.T) {
	dto := ClusterDTORequest{
		ClusterName: "", // NewCluster will fail
		InfraID:     "ABCDE",
		Provider:    inventory.AWSProvider,
		Region:      "eu-west-1",
		ConsoleLink: "https://console",
		Owner:       "owner",
	}

	cluster := dto.ToInventoryCluster()
	assert.Nil(t, cluster)
}

// TestToInventoryClusterList verifies slice conversion from DTOs to inventory.Cluster.
func TestToInventoryClusterList(t *testing.T) {
	t.Run("Multiple DTOs", func(t *testing.T) { testToInventoryClusterList_Correct(t) })
}

func testToInventoryClusterList_Correct(t *testing.T) {
	now := time.Now().UTC()

	dtos := []ClusterDTORequest{
		{
			ClusterName: "c1",
			InfraID:     "AAAAA",
			Provider:    inventory.AWSProvider,
			Status:      inventory.Running,
			Region:      "eu-west-1",
			AccountID:   "acc-1",
			ConsoleLink: "https://console-1",
			LastScanTS:  now,
			CreatedAt:   now.Add(-time.Hour),
			Owner:       "owner-1",
		},
		{
			ClusterName: "c2",
			InfraID:     "BBBBB",
			Provider:    inventory.AWSProvider,
			Status:      inventory.Stopped,
			Region:      "us-east-1",
			AccountID:   "acc-2",
			ConsoleLink: "https://console-2",
			LastScanTS:  now.Add(-2 * time.Hour),
			CreatedAt:   now.Add(-3 * time.Hour),
			Owner:       "owner-2",
		},
	}

	clusters := ToInventoryClusterList(dtos)

	assert.NotNil(t, clusters)
	assert.Len(t, *clusters, 2)

	assert.Equal(t, "c1", (*clusters)[0].ClusterName)
	assert.Equal(t, "c2", (*clusters)[1].ClusterName)
}

// TestToClusterDTORequest verifies inventory.Cluster to DTO conversion.
func TestToClusterDTORequest(t *testing.T) {
	t.Run("Cluster to DTO", func(t *testing.T) { testToClusterDTORequest_Correct(t) })
}

func testToClusterDTORequest_Correct(t *testing.T) {
	now := time.Now().UTC()

	cluster := inventory.Cluster{
		ClusterID:   "cluster-1",
		ClusterName: "name-1",
		InfraID:     "ABCDE",
		Provider:    inventory.AWSProvider,
		Status:      inventory.Running,
		Region:      "eu-west-1",
		AccountID:   "acc-1",
		ConsoleLink: "https://console",
		LastScanTS:  now,
		CreatedAt:   now.Add(-time.Hour),
		Age:         77,
		Owner:       "owner",
	}

	dto := ToClusterDTORequest(cluster)

	assert.NotNil(t, dto)
	assert.Equal(t, cluster.ClusterID, dto.ClusterID)
	assert.Equal(t, cluster.ClusterName, dto.ClusterName)
	assert.Equal(t, cluster.InfraID, dto.InfraID)
	assert.Equal(t, cluster.Provider, dto.Provider)
	assert.Equal(t, cluster.Status, dto.Status)
	assert.Equal(t, cluster.Region, dto.Region)
	assert.Equal(t, cluster.AccountID, dto.AccountID)
	assert.Equal(t, cluster.ConsoleLink, dto.ConsoleLink)
	assert.Equal(t, cluster.LastScanTS, dto.LastScanTS)
	assert.Equal(t, cluster.CreatedAt, dto.CreatedAt)
	assert.Equal(t, cluster.Age, dto.Age)
	assert.Equal(t, cluster.Owner, dto.Owner)
}

// TestToClusterDTORequestList verifies list conversion from inventory clusters to DTO requests.
func TestToClusterDTORequestList(t *testing.T) {
	t.Run("Cluster DTO list", func(t *testing.T) { testToClusterDTORequestList_Correct(t) })
}

func testToClusterDTORequestList_Correct(t *testing.T) {
	clusters := []inventory.Cluster{
		{ClusterID: "c1", ClusterName: "name-1", InfraID: "AAAAA", Provider: inventory.AWSProvider},
		{ClusterID: "c2", ClusterName: "name-2", InfraID: "BBBBB", Provider: inventory.AWSProvider},
	}

	dtoList := ToClusterDTORequestList(clusters)

	assert.NotNil(t, dtoList)
	assert.Len(t, *dtoList, 2)
	assert.Equal(t, "c1", (*dtoList)[0].ClusterID)
	assert.Equal(t, "c2", (*dtoList)[1].ClusterID)
}
