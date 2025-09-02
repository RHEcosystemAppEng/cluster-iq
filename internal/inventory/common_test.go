package inventory

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	t.Run("ZeroDays", func(t *testing.T) { testCalculateAge_ZeroDays(t) })
	t.Run("Normal", func(t *testing.T) { testCalculateAge_Normal(t) })
	t.Run("WithinOneDay", func(t *testing.T) { testCalculateAge_WithinOneDay(t) })
}

// testCalculateAge_ZeroDays verifies age is 1 when timestamps are equal (zero duration).
func testCalculateAge_ZeroDays(t *testing.T) {
	now := time.Now()
	age := calculateAge(now, now)
	if age != 1 {
		t.Errorf("expected age 1 for same timestamps, got %d", age)
	}
}

// testCalculateAge_Normal verifies age is correct when time difference is > 24h.
func testCalculateAge_Normal(t *testing.T) {
	init := time.Now().Add(-48 * time.Hour)
	end := time.Now()
	age := calculateAge(init, end)
	if age != 2 {
		t.Errorf("expected age 2, got %d", age)
	}
}

// testCalculateAge_WithinOneDay verifies that a diff < 24h still returns age 1.
func testCalculateAge_WithinOneDay(t *testing.T) {
	init := time.Now().Add(-4 * time.Hour)
	end := time.Now()
	age := calculateAge(init, end)
	if age != 1 {
		t.Errorf("expected age 1 for <24h difference, got %d", age)
	}
}
