package db_test

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/convert"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/stretchr/testify/assert"
)

// TestAccountDBResponse_ToAccountDTOResponse verifies DB model to DTO conversion.
func TestAccountDBResponse_ToAccountDTOResponse(t *testing.T) {
	t.Run("Convert AccountDBResponse to AccountDTOResponse", func(t *testing.T) {
		testAccountDBResponse_ToAccountDTOResponse_Correct(t)
	})
}

func testAccountDBResponse_ToAccountDTOResponse_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	model := db.AccountDBResponse{
		AccountID:             "acc-1",
		AccountName:           "test-account",
		Provider:              inventory.AWSProvider,
		LastScanTimestamp:     now,
		CreatedAt:             now.Add(-time.Hour),
		ClusterCount:          3,
		TotalCost:             123.45,
		Last15DaysCost:        45.67,
		LastMonthCost:         89.01,
		CurrentMonthSoFarCost: 12.34,
	}

	dto := conv.ToAccountDTO(model)

	assert.Equal(t, model.AccountID, dto.AccountID)
	assert.Equal(t, model.AccountName, dto.AccountName)
	assert.Equal(t, model.Provider, dto.Provider)
	assert.Equal(t, model.LastScanTimestamp, dto.LastScanTimestamp)
	assert.Equal(t, model.CreatedAt, dto.CreatedAt)
	assert.Equal(t, model.ClusterCount, dto.ClusterCount)
	assert.Equal(t, model.TotalCost, dto.TotalCost)
	assert.Equal(t, model.Last15DaysCost, dto.Last15DaysCost)
	assert.Equal(t, model.LastMonthCost, dto.LastMonthCost)
	assert.Equal(t, model.CurrentMonthSoFarCost, dto.CurrentMonthSoFarCost)
}

// TestToAccountDTOResponseList verifies slice conversion from DB responses to DTOs.
func TestToAccountDTOResponseList(t *testing.T) {
	t.Run("Convert AccountDBResponse list to DTO list", func(t *testing.T) {
		testToAccountDTOResponseList_Correct(t)
	})
}

func testToAccountDTOResponseList_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	models := []db.AccountDBResponse{
		{
			AccountID:         "acc-1",
			AccountName:       "account-1",
			Provider:          inventory.AWSProvider,
			LastScanTimestamp: now,
			CreatedAt:         now.Add(-time.Hour),
			ClusterCount:      1,
		},
		{
			AccountID:         "acc-2",
			AccountName:       "account-2",
			Provider:          inventory.AWSProvider,
			LastScanTimestamp: now.Add(-time.Minute),
			CreatedAt:         now.Add(-2 * time.Hour),
			ClusterCount:      2,
		},
	}

	dtos := conv.ToAccountDTOs(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, "acc-1", dtos[0].AccountID)
	assert.Equal(t, "acc-2", dtos[1].AccountID)
	assert.Equal(t, 1, dtos[0].ClusterCount)
	assert.Equal(t, 2, dtos[1].ClusterCount)
}
