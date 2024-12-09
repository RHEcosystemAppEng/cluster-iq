package inventory

import (
	"testing"
	"time"
)

func TestNewExpense(t *testing.T) {
	var tests = []struct {
		id         string
		cost       float64
		date       time.Time
		returnsNil bool
	}{
		{
			"01234",
			12.5,
			time.Now(),
			false,
		},
		{
			"01234",
			-12.5,
			time.Now(),
			true,
		},
	}

	for _, test := range tests {
		expense := NewExpense(test.id, test.cost, test.date)
		if (expense == nil) != test.returnsNil {
			t.Errorf("Expense created is null")
			return
		}
		if expense != nil && expense.Amount < 0 {
			t.Errorf("Expense's Cost is below zero. Have: %f ; Expected: %f", expense.Amount, test.cost)
		}
	}
}
