package inventory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExpense(t *testing.T) {
	t.Run("Normal", func(t *testing.T) { testNewExpenseNormal(t) })
	t.Run("NegativeAmount", func(t *testing.T) { testNewExpenseNegativeAmount(t) })
}

// testNewExpenseNormal verifies the expense is created correctly in normal conditions
func testNewExpenseNormal(t *testing.T) {
	amount := 12.5
	expense := NewExpense("instance", amount, time.Now())

	assert.NotNil(t, expense)
	assert.Equal(t, expense.Amount, 12.5)
}

// testNewExpenseNegativeAmount verifies the expense doesn't accept negative amounts
func testNewExpenseNegativeAmount(t *testing.T) {
	amount := -12.5
	expense := NewExpense("instance", amount, time.Now())

	assert.NotNil(t, expense)
	assert.Equal(t, expense.Amount, 0.0)
}
