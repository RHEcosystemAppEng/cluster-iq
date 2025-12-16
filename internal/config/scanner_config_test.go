package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLoadScannerConfig verifies loading ScannerConfig from environment variables.
func TestLoadScannerConfig(t *testing.T) {
	t.Run("Load ScannerConfig OK", func(t *testing.T) { testLoadScannerConfig_OK(t) })
	t.Run("Load ScannerConfig defaults", func(t *testing.T) { testLoadScannerConfig_Defaults(t) })
	t.Run("Load ScannerConfig missing required var", func(t *testing.T) { testLoadScannerConfig_MissingVar(t) })
}

func testLoadScannerConfig_OK(t *testing.T) {
	t.Setenv("CIQ_API_URL", "https://api.clusteriq.local")
	t.Setenv("CIQ_SKIP_NO_OPENSHIFT_INSTANCES", "false")
	t.Setenv("CIQ_CREDS_FILE", "./test-file")

	cfg, err := LoadScannerConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "https://api.clusteriq.local", cfg.APIURL)
	assert.False(t, cfg.SkipNoOpenShiftInstances)
}

func testLoadScannerConfig_Defaults(t *testing.T) {
	clearScannerEnv()

	// Required
	t.Setenv("CIQ_API_URL", "https://api.clusteriq.local")
	t.Setenv("CIQ_CREDS_FILE", "./test-file")

	// Not set:
	// - CIQ_SKIP_NO_OPENSHIFT_INSTANCES (default true)

	cfg, err := LoadScannerConfig()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "https://api.clusteriq.local", cfg.APIURL)
	assert.True(t, cfg.SkipNoOpenShiftInstances)
}

func testLoadScannerConfig_MissingVar(t *testing.T) {
	clearScannerEnv()

	cfg, err := LoadScannerConfig()

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func clearScannerEnv() {
	_ = os.Unsetenv("CIQ_API_URL")
	_ = os.Unsetenv("CIQ_SKIP_NO_OPENSHIFT_INSTANCES")
	_ = os.Unsetenv("CIQ_CREDS_FILE")
}
