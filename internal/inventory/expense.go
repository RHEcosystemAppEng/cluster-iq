package inventory

import "time"

// Expense defines the expenses applied to an instance
type Expense struct {
	// InstanceID references the instance of the expense
	InstanceID string `db:"instance_id" json:"instanceID"`

	// Ammount represents the cost in USDollars
	Amount float64 `db:"amount" json:"amount"`

	// Date (Year, month, day)
	Date time.Time `db:"date" json:"date"`
}

// NewExpense create a expense for an instance
func NewExpense(instanceID string, amount float64, date time.Time) *Expense {
	// Checking if cost is below zero, which is not possible
	if amount < 0.0 {
		amount = 0.0
	}

	return &Expense{
		InstanceID: instanceID,
		Amount:     amount,
		Date:       date,
	}
}
