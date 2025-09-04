package inventory

import "time"

// OverviewSummary represents the comprehensive overview of the system's inventory.
type OverviewSummary struct {
	Clusters  ClustersSummary
	Instances InstancesSummary
	Providers ProvidersSummary
	Scanner   Scanner
}

// ClustersSummary provides a summary of cluster counts by status.
type ClustersSummary struct {
	Running  int `db:"running"`
	Stopped  int `db:"stopped"`
	Archived int `db:"archived"`
}

// InstancesSummary provides a summary of instance counts.
type InstancesSummary struct {
	Running  int `db:"running"`
	Stopped  int `db:"stopped"`
	Archived int `db:"archived"`
}

// ProvidersSummary provides a summary of provider-specific data.
type ProvidersSummary struct {
	AWS   ProviderDetails
	GCP   ProviderDetails
	Azure ProviderDetails
}

// ProviderDetails contains the account and cluster counts for a specific provider.
type ProviderDetails struct {
	AccountCount int
	ClusterCount int
}

// Scanner provides information about the last inventory scan.
type Scanner struct {
	LastScanTimestamp time.Time
}
