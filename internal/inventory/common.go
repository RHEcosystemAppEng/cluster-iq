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

// Calculates Age parameter in days
func calculateAge(initTimestamp time.Time, finalTimestamp time.Time) int {
	return int(finalTimestamp.Sub(initTimestamp).Hours() / 24)
}
