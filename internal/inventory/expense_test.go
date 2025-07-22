package inventory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewExpense(t *testing.T) {
	instanceID := "testInstance"
	amount := 12.34
	date := time.Now()

	expectedExpense := &Expense{
		InstanceID: instanceID,
		Amount:     amount,
		Date:       date,
	}

	actualExpense := NewExpense(instanceID, amount, date)

	assert.NotNil(t, actualExpense)
	assert.Equal(t, expectedExpense, actualExpense)

	wrongExpense := NewExpense(instanceID, -6.0, date)
	assert.Nil(t, wrongExpense)
}
