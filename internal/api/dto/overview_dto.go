package dto

import "time"

// OverviewSummary represents the comprehensive overview of the system's inventory.
type OverviewSummary struct {
	Clusters  ClusterSummary   `json:"clusters"`
	Instances InstancesSummary `json:"instances"`
	Providers ProvidersSummary `json:"providers"`
	Scanner   Scanner          `json:"scanner"`
}

// ClusterSummary provides a summary of cluster counts by status.
type ClusterSummary struct {
	Running  int `json:"running"`
	Stopped  int `json:"stopped"`
	Archived int `json:"archived"`
}

// InstancesSummary provides a summary of instance counts by status.
type InstancesSummary struct {
	Running  int `json:"running"`
	Stopped  int `json:"stopped"`
	Archived int `json:"archived"`
}

// ProvidersSummary provides a summary of provider-specific data.
type ProvidersSummary struct {
	AWS   ProviderDetails `json:"aws"`
	GCP   ProviderDetails `json:"gcp"`
	Azure ProviderDetails `json:"azure"`
}

// ProviderDetails contains the account and cluster counts for a specific provider.
type ProviderDetails struct {
	AccountCount int `json:"account_count"`
	ClusterCount int `json:"cluster_count"`
}

// Scanner provides information about the last inventory scan.
type Scanner struct {
	LastScanTimestamp time.Time `json:"last_scan_timestamp"`
}
