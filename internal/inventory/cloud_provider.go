package inventory

import "strings"

// CloudProvider defines the cloud provider of the instance
type CloudProvider string

const (
	// AWSProvider - Amazon Web Services Cloud Provider
	AWSProvider CloudProvider = "AWS"
	// AzureProvider - Microsoft Azure Cloud Provider
	AzureProvider = "Azure"
	// GCPProvider - Google Cloud Platform Cloud Provider
	GCPProvider = "GCP"
	// UnknownProvider - Unknown Platform Cloud Provider
	UnknownProvider = "UNKNOWN"
)

// GetCloudProvider checks a incoming string and returns the corresponding inventory.CloudProvider value
func GetCloudProvider(provider string) CloudProvider {
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
