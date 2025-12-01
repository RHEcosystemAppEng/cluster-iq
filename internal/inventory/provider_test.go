package inventory

import (
	"testing"
)

// TODO REWRITE NEW FORMAT
func TestGetCloudProvider(t *testing.T) {
	tests := []struct {
		input    string
		expected Provider
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
			result := GetProvider(tt.input)
			if result != tt.expected {
				t.Errorf("GetProvider(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
