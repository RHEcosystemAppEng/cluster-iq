package inventory

import "strings"

// CloudProvider defines the supported cloud providers on the inventory
type CloudProvider string

const (
	AWSProvider     CloudProvider = "AWS"     // AWSProvider - Amazon Web Services Cloud Provider
	AzureProvider   CloudProvider = "Azure"   // AzureProvider - Microsoft Azure Cloud Provider
	GCPProvider     CloudProvider = "GCP"     // GCPProvider - Google Cloud Platform Cloud Provider
	UnknownProvider CloudProvider = "UNKNOWN" // UnknownProvider - Unknown Platform Cloud Provider
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
