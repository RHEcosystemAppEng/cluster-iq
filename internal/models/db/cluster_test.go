package db_test

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/convert"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/stretchr/testify/assert"
)

// TestClusterDBResponse_ToClusterDTOResponse verifies DB model to DTO conversion.
func TestClusterDBResponse_ToClusterDTOResponse(t *testing.T) {
	t.Run("Convert ClusterDBResponse to ClusterDTOResponse", func(t *testing.T) {
		testClusterDBResponse_ToClusterDTOResponse_Correct(t)
	})
}

func testClusterDBResponse_ToClusterDTOResponse_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	model := db.ClusterDBResponse{
		ClusterID:             "cluster-1",
		ClusterName:           "test-cluster",
		InfraID:               "ABCDE",
		Provider:              inventory.AWSProvider,
		Status:                inventory.Running,
		Region:                "eu-west-1",
		AccountID:             "acc-1",
		AccountName:           "account-1",
		ConsoleLink:           "https://console",
		LastScanTimestamp:     now,
		CreatedAt:             now.Add(-time.Hour),
		Age:                   42,
		Owner:                 "team-a",
		InstanceCount:         5,
		TotalCost:             123.45,
		Last15DaysCost:        12.34,
		LastMonthCost:         56.78,
		CurrentMonthSoFarCost: 9.01,
	}

	dto := conv.ToClusterDTO(model)

	assert.Equal(t, model.ClusterID, dto.ClusterID)
	assert.Equal(t, model.ClusterName, dto.ClusterName)
	assert.Equal(t, model.InfraID, dto.InfraID)
	assert.Equal(t, model.Provider, dto.Provider)
	assert.Equal(t, model.Status, dto.Status)
	assert.Equal(t, model.Region, dto.Region)
	assert.Equal(t, model.AccountID, dto.AccountID)
	assert.Equal(t, model.ConsoleLink, dto.ConsoleLink)
	assert.Equal(t, model.InstanceCount, dto.InstanceCount)
	assert.Equal(t, model.LastScanTimestamp, dto.LastScanTimestamp)
	assert.Equal(t, model.CreatedAt, dto.CreatedAt)
	assert.Equal(t, model.Age, dto.Age)
	assert.Equal(t, model.Owner, dto.Owner)
	assert.Equal(t, model.TotalCost, dto.TotalCost)
	assert.Equal(t, model.Last15DaysCost, dto.Last15DaysCost)
	assert.Equal(t, model.LastMonthCost, dto.LastMonthCost)
	assert.Equal(t, model.CurrentMonthSoFarCost, dto.CurrentMonthSoFarCost)
}

// TestToClusterDTOResponseList verifies slice conversion from DB responses to DTOs.
func TestToClusterDTOResponseList(t *testing.T) {
	t.Run("Convert ClusterDBResponse list to DTO list", func(t *testing.T) {
		testToClusterDTOResponseList_Correct(t)
	})
}

func testToClusterDTOResponseList_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	models := []db.ClusterDBResponse{
		{
			ClusterID:         "c1",
			ClusterName:       "cluster-1",
			InfraID:           "AAAAA",
			Provider:          inventory.AWSProvider,
			Status:            inventory.Running,
			Region:            "eu-west-1",
			AccountID:         "acc-1",
			ConsoleLink:       "https://console-1",
			LastScanTimestamp: now,
			CreatedAt:         now.Add(-time.Hour),
			InstanceCount:     1,
		},
		{
			ClusterID:         "c2",
			ClusterName:       "cluster-2",
			InfraID:           "BBBBB",
			Provider:          inventory.AWSProvider,
			Status:            inventory.Stopped,
			Region:            "us-east-1",
			AccountID:         "acc-2",
			ConsoleLink:       "https://console-2",
			LastScanTimestamp: now.Add(-time.Minute),
			CreatedAt:         now.Add(-2 * time.Hour),
			InstanceCount:     2,
		},
	}

	dtos := conv.ToClusterDTOs(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, "c1", dtos[0].ClusterID)
	assert.Equal(t, "c2", dtos[1].ClusterID)
	assert.Equal(t, 1, dtos[0].InstanceCount)
	assert.Equal(t, 2, dtos[1].InstanceCount)
}
