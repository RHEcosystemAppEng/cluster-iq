package inventory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAge(t *testing.T) {
	t.Run("ZeroDays", func(t *testing.T) { testCalculateAge_ZeroDays(t) })
	t.Run("Normal", func(t *testing.T) { testCalculateAge_Normal(t) })
}

// testCalculateAge_ZeroDays verifies age is 1 when timestamps are equal (zero duration).
func testCalculateAge_ZeroDays(t *testing.T) {
	now := time.Now()
	age := calculateAge(now, now)

	assert.Equal(t, age, 1)
}

// testCalculateAge_Normal verifies age is correct when time difference is > 24h.
func testCalculateAge_Normal(t *testing.T) {
	init := time.Now().Add(-48 * time.Hour)
	end := time.Now()
	age := calculateAge(init, end)

	assert.Equal(t, age, 2)
}
