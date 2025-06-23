package inventory

import (
	"encoding/json"
	"time"
)

// JSONMarshal converts an inventory object as JSON format for printing
func JSONMarshal(object interface{}) (string, error) {
	b, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Calculates Age parameter in days If the resulting age is 0, 1 is returned
// instead. This is because when calculating daily_cost, if the instance has age
// 0, it will throw a 'divide by zero' error
func calculateAge(initTimestamp time.Time, finalTimestamp time.Time) int {
	age := int(finalTimestamp.Sub(initTimestamp).Hours() / 24)
	if age == 0 {
		return 1
	}
	return age
}
