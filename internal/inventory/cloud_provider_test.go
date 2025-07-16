package inventory

import (
	"testing"
)

func TestGetCloudProvider(t *testing.T) {
	tests := []struct {
		input    string
		expected CloudProvider
	}{
		{"aws", AWSProvider},
		{"AWS", AWSProvider},
		{"gcp", GCPProvider},
		{"GCP", GCPProvider},
		{"azure", AzureProvider},
		{"AZURE", AzureProvider},
		{"unknown", UnknownProvider},
		{"", UnknownProvider},
		{"DigitalOcean", UnknownProvider},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := GetCloudProvider(tt.input)
			if result != tt.expected {
				t.Errorf("GetCloudProvider(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
