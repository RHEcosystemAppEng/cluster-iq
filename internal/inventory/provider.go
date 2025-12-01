package inventory

import "strings"

// Provider defines the different values for the supported cloud/infrastructure providers
type Provider string

const (
	// AWSProvider - Amazon Web Services Cloud Provider
	AWSProvider Provider = "AWS"

	// AzureProvider - Microsoft Azure Cloud Provider
	AzureProvider Provider = "Azure"

	// GCPProvider - Google Cloud Platform Cloud Provider
	GCPProvider Provider = "GCP"

	// UnknownProvider - Unknown Provider
	UnknownProvider Provider = "UNKNOWN_PROVIDER"
)

// GetProvider checks a incoming string and returns the corresponding inventory.Provider value
func GetProvider(provider string) Provider {
	switch strings.ToUpper(provider) {
	case "AWS":
		return AWSProvider
	case "GCP":
		return GCPProvider
	case "AZURE":
		return AzureProvider
	default:
		return UnknownProvider
	}
}
