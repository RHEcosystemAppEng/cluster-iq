package inventory

import (
	"time"
)

const (
	HoursPerDay = 24
)

// Calculates Age parameter in days If the resulting age is 0, 1 is returned
// instead. This is because when calculating daily_cost, if the instance has age
// 0, it will throw a 'divide by zero' error
func calculateAge(initTimestamp time.Time, finalTimestamp time.Time) int {
	age := int(finalTimestamp.Sub(initTimestamp).Hours() / HoursPerDay)
	if age == 0 {
		return 1
	}

	return age
}
