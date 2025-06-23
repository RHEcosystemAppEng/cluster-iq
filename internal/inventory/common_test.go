package inventory

import (
	"strings"
	"testing"
	"time"
)

// TestJSONMarshal_Success verifies that JSONMarshal correctly serializes an object to indented JSON.
func TestJSONMarshal_Success(t *testing.T) {
	obj := struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{
		Name:  "test",
		Value: 42,
	}

	jsonStr, err := JSONMarshal(obj)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if jsonStr == "" {
		t.Error("expected non-empty JSON string")
	}
	if !(containsAll(jsonStr, "{", "\"name\"", "\"value\"", "test", "42")) {
		t.Errorf("unexpected JSON output: %s", jsonStr)
	}
}

// helper to verify multiple substrings in a string
func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// TestJSONMarshal_Error verifies that JSONMarshal returns an error when marshalling fails.
func TestJSONMarshal_Error(t *testing.T) {
	// Channels can't be marshalled to JSON
	ch := make(chan int)
	_, err := JSONMarshal(ch)
	if err == nil {
		t.Error("expected error when marshalling unsupported type, got nil")
	}
}

// TestCalculateAge_ZeroDays verifies that age is 1 when timestamps are equal (zero duration).
func TestCalculateAge_ZeroDays(t *testing.T) {
	now := time.Now()
	age := calculateAge(now, now)
	if age != 1 {
		t.Errorf("expected age 1 for same timestamps, got %d", age)
	}
}

// TestCalculateAge_Normal verifies that age is correct when time difference is > 24h.
func TestCalculateAge_Normal(t *testing.T) {
	init := time.Now().Add(-48 * time.Hour)
	end := time.Now()
	age := calculateAge(init, end)
	if age != 2 {
		t.Errorf("expected age 2, got %d", age)
	}
}

// TestCalculateAge_WithinOneDay verifies that a diff < 24h still returns age 1.
func TestCalculateAge_WithinOneDay(t *testing.T) {
	init := time.Now().Add(-4 * time.Hour)
	end := time.Now()
	age := calculateAge(init, end)
	if age != 1 {
		t.Errorf("expected age 1 for <24h difference, got %d", age)
	}
}
