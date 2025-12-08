package dto

import (
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// OverviewSummary represents the comprehensive overview of the system's inventory.
type OverviewSummary struct {
	Clusters  ClusterSummary   `json:"clusters"`
	Instances InstancesSummary `json:"instances"`
	Providers ProvidersSummary `json:"providers"`
	Scanner   Scanner          `json:"scanner"`
} // @name OverviewSummary

// ClusterSummary provides a summary of cluster counts by status.
type ClusterSummary struct {
	Running  int `json:"running"`
	Stopped  int `json:"stopped"`
	Archived int `json:"archived"`
} // @name ClusterSummary

// InstancesSummary provides a summary of instance counts by status.
type InstancesSummary struct {
	Running  int `json:"running"`
	Stopped  int `json:"stopped"`
	Archived int `json:"archived"`
} // @name InstancesSummary

// ProvidersSummary provides a summary of provider-specific data.
type ProvidersSummary struct {
	AWS   ProviderDetails `json:"aws"`
	GCP   ProviderDetails `json:"gcp"`
	Azure ProviderDetails `json:"azure"`
} // @name ProvidersSummary

// ProviderDetails contains the account and cluster counts for a specific provider.
type ProviderDetails struct {
	AccountCount int `json:"account_count"`
	ClusterCount int `json:"cluster_count"`
} // @name ProviderDetails

// Scanner provides information about the last inventory scan.
type Scanner struct {
	LastScanTimestamp time.Time `json:"lastScanTimestamp"`
} // @name Scanner

// ToOverviewSummaryDTO converts an inventory OverviewSummary to a DTO.
func ToOverviewSummaryDTO(model inventory.OverviewSummary) OverviewSummary {
	return OverviewSummary{
		Clusters:  toClusterSummaryDTO(model.Clusters),
		Instances: toInstancesSummaryDTO(model.Instances),
		Providers: toProvidersSummaryDTO(model.Providers),
		Scanner:   toScannerDTO(model.Scanner),
	}
}

func toClusterSummaryDTO(model inventory.ClustersSummary) ClusterSummary {
	return ClusterSummary{
		Running:  model.Running,
		Stopped:  model.Stopped,
		Archived: model.Archived,
	}
}

func toInstancesSummaryDTO(model inventory.InstancesSummary) InstancesSummary {
	return InstancesSummary{
		Running:  model.Running,
		Stopped:  model.Stopped,
		Archived: model.Archived,
	}
}

func toProvidersSummaryDTO(model inventory.ProvidersSummary) ProvidersSummary {
	return ProvidersSummary{
		AWS:   toProviderDetailsDTO(model.AWS),
		GCP:   toProviderDetailsDTO(model.GCP),
		Azure: toProviderDetailsDTO(model.Azure),
	}
}

func toProviderDetailsDTO(model inventory.ProviderDetails) ProviderDetails {
	return ProviderDetails{
		AccountCount: model.AccountCount,
		ClusterCount: model.ClusterCount,
	}
}

func toScannerDTO(model inventory.Scanner) Scanner {
	return Scanner{
		LastScanTimestamp: model.LastScanTimestamp,
	}
}
