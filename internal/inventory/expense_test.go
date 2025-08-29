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
	expense := NewExpense("instance", 12.5, time.Now())
	assert.NotNil(t, expense)
	assert.GreaterOrEqual(t, expense.Amount, 0.0)
}

// testNewExpenseNegativeAmount verifies the expense doesn't accept negative amounts
func testNewExpenseNegativeAmount(t *testing.T) {
	expense := NewExpense("instance", -12.5, time.Now())
	assert.NotNil(t, expense)
	assert.Equal(t, expense.Amount, 0.0)
}
