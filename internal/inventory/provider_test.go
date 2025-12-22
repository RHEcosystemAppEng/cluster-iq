package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCloudProvider(t *testing.T) {
	t.Run("AWS", func(t *testing.T) { testGetCloudProvider_AWS(t) })
	t.Run("Azure", func(t *testing.T) { testGetCloudProvider_Azure(t) })
	t.Run("GCP", func(t *testing.T) { testGetCloudProvider_GCP(t) })
	t.Run("Unknown", func(t *testing.T) { testGetCloudProvider_Unknown(t) })
}

func testGetCloudProvider_AWS(t *testing.T) {
	assert.Equal(t, AWSProvider, GetProvider("AWS"))
	assert.NotEqual(t, AWSProvider, GetProvider("Azure"))
	assert.Equal(t, AWSProvider, GetProvider("AwS"))
	assert.NotEqual(t, AWSProvider, GetProvider(""))
}

func testGetCloudProvider_Azure(t *testing.T) {
	assert.Equal(t, AzureProvider, GetProvider("Azure"))
	assert.NotEqual(t, AzureProvider, GetProvider("Google"))
	assert.Equal(t, AzureProvider, GetProvider("AZurE"))
	assert.NotEqual(t, AzureProvider, GetProvider(""))
}

func testGetCloudProvider_GCP(t *testing.T) {
	assert.Equal(t, GCPProvider, GetProvider("GCP"))
	assert.NotEqual(t, GCPProvider, GetProvider("Azure"))
	assert.Equal(t, GCPProvider, GetProvider("gCp"))
	assert.NotEqual(t, GCPProvider, GetProvider(""))
}

func testGetCloudProvider_Unknown(t *testing.T) {
	assert.Equal(t, UnknownProvider, GetProvider("no-idea"))
	assert.Equal(t, UnknownProvider, GetProvider("DigitalOcean"))
	assert.Equal(t, UnknownProvider, GetProvider(""))
}
