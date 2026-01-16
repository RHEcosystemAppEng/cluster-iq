package db_test

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/convert"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/stretchr/testify/assert"
)

// TestExpenseDBResponse_ToExpenseDTOResponse verifies DB model to DTO conversion.
func TestExpenseDBResponse_ToExpenseDTOResponse(t *testing.T) {
	t.Run("Convert ExpenseDBResponse to ExpenseDTOResponse", func(t *testing.T) {
		testExpenseDBResponse_ToExpenseDTOResponse_Correct(t)
	})
}

func testExpenseDBResponse_ToExpenseDTOResponse_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	model := db.ExpenseDBResponse{
		InstanceID: "i-123",
		Amount:     42.5,
		Date:       now,
	}

	dto := conv.ToExpenseDTO(model)

	assert.Equal(t, model.InstanceID, dto.InstanceID)
	assert.Equal(t, model.Amount, dto.Amount)
	assert.Equal(t, model.Date, dto.Date)
}

// TestToExpenseDTOResponseList verifies slice conversion from DB responses to DTOs.
func TestToExpenseDTOResponseList(t *testing.T) {
	t.Run("Convert ExpenseDBResponse list to DTO list", func(t *testing.T) {
		testToExpenseDTOResponseList_Correct(t)
	})
}

func testToExpenseDTOResponseList_Correct(t *testing.T) {
	now := time.Now().UTC()
	conv := &convert.ConverterImpl{}

	models := []db.ExpenseDBResponse{
		{
			InstanceID: "i-1",
			Amount:     10.0,
			Date:       now,
		},
		{
			InstanceID: "i-2",
			Amount:     20.5,
			Date:       now.Add(-24 * time.Hour),
		},
	}

	dtos := conv.ToExpenseDTOs(models)

	assert.Len(t, dtos, 2)
	assert.Equal(t, "i-1", dtos[0].InstanceID)
	assert.Equal(t, 10.0, dtos[0].Amount)
	assert.Equal(t, "i-2", dtos[1].InstanceID)
	assert.Equal(t, 20.5, dtos[1].Amount)
}
