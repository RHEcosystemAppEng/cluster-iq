package dto

import (
	"testing"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// TestExpenseDTORequest_ToInventoryExpense verifies DTO to inventory.Expense conversion.
func TestExpenseDTORequest_ToInventoryExpense(t *testing.T) {
	t.Run("Valid DTO", func(t *testing.T) { testExpenseDTORequest_ToInventoryExpense_Correct(t) })
	t.Run("Invalid DTO returns nil", func(t *testing.T) { testExpenseDTORequest_ToInventoryExpense_Invalid(t) })
}

func testExpenseDTORequest_ToInventoryExpense_Correct(t *testing.T) {
	now := time.Now().UTC()

	dto := ExpenseDTORequest{
		InstanceID: "i-1",
		Amount:     12.5,
		Date:       now,
	}

	exp := dto.ToInventoryExpense()

	assert.NotNil(t, exp)
	assert.Equal(t, dto.InstanceID, exp.InstanceID)
	assert.Equal(t, dto.Amount, exp.Amount)
	assert.Equal(t, dto.Date, exp.Date)
}

func testExpenseDTORequest_ToInventoryExpense_Invalid(t *testing.T) {
	dto := ExpenseDTORequest{
		InstanceID: "i-1",
		Amount:     -10.0, // NewExpense should reject negative amounts
		Date:       time.Now().UTC(),
	}

	exp := dto.ToInventoryExpense()
	assert.NotNil(t, exp)
	assert.Zero(t, exp.Amount)
}

// TestToInventoryExpenseList verifies slice conversion from DTOs to inventory.Expense.
func TestToInventoryExpenseList(t *testing.T) {
	t.Run("Multiple DTOs", func(t *testing.T) { testToInventoryExpenseList_Correct(t) })
}

func testToInventoryExpenseList_Correct(t *testing.T) {
	now := time.Now().UTC()

	dtos := []ExpenseDTORequest{
		{InstanceID: "i-1", Amount: 1.2, Date: now},
		{InstanceID: "i-2", Amount: 3.4, Date: now.Add(-time.Hour)},
	}

	expenses := ToInventoryExpenseList(dtos)

	assert.NotNil(t, expenses)
	assert.Len(t, *expenses, 2)
	assert.Equal(t, "i-1", (*expenses)[0].InstanceID)
	assert.Equal(t, "i-2", (*expenses)[1].InstanceID)
}

// TestToExpenseDTORequest verifies inventory.Expense to DTO conversion.
func TestToExpenseDTORequest(t *testing.T) {
	t.Run("Expense to DTO", func(t *testing.T) { testToExpenseDTORequest_Correct(t) })
}

func testToExpenseDTORequest_Correct(t *testing.T) {
	now := time.Now().UTC()

	exp := inventory.Expense{
		InstanceID: "i-1",
		Amount:     9.99,
		Date:       now,
	}

	dto := ToExpenseDTORequest(exp)

	assert.NotNil(t, dto)
	assert.Equal(t, exp.InstanceID, dto.InstanceID)
	assert.Equal(t, exp.Amount, dto.Amount)
	assert.Equal(t, exp.Date, dto.Date)
}

// TestToExpenseDTORequestList verifies list conversion from inventory expenses to DTOs.
func TestToExpenseDTORequestList(t *testing.T) {
	t.Run("Expense DTO list", func(t *testing.T) { testToExpenseDTORequestList_Correct(t) })
}

func testToExpenseDTORequestList_Correct(t *testing.T) {
	expenses := []inventory.Expense{
		{InstanceID: "i-1", Amount: 1.0, Date: time.Now()},
		{InstanceID: "i-2", Amount: 2.0, Date: time.Now()},
	}

	dtoList := ToExpenseDTORequestList(expenses)

	assert.NotNil(t, dtoList)
	assert.Len(t, *dtoList, 2)
	assert.Equal(t, "i-1", (*dtoList)[0].InstanceID)
	assert.Equal(t, "i-2", (*dtoList)[1].InstanceID)
}
