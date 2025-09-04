package mappers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
)

// ToOverviewSummaryDTO converts an inventory OverviewSummary to a DTO.
func ToOverviewSummaryDTO(model inventory.OverviewSummary) dto.OverviewSummary {
	return dto.OverviewSummary{
		Clusters:  toClusterSummaryDTO(model.Clusters),
		Instances: toInstancesSummaryDTO(model.Instances),
		Providers: toProvidersSummaryDTO(model.Providers),
		Scanner:   toScannerDTO(model.Scanner),
	}
}

func toClusterSummaryDTO(model inventory.ClustersSummary) dto.ClusterSummary {
	return dto.ClusterSummary{
		Running:  model.Running,
		Stopped:  model.Stopped,
		Archived: model.Archived,
	}
}

func toInstancesSummaryDTO(model inventory.InstancesSummary) dto.InstancesSummary {
	return dto.InstancesSummary{
		Running:  model.Running,
		Stopped:  model.Stopped,
		Archived: model.Archived,
	}
}

func toProvidersSummaryDTO(model inventory.ProvidersSummary) dto.ProvidersSummary {
	return dto.ProvidersSummary{
		AWS:   toProviderDetailsDTO(model.AWS),
		GCP:   toProviderDetailsDTO(model.GCP),
		Azure: toProviderDetailsDTO(model.Azure),
	}
}

func toProviderDetailsDTO(model inventory.ProviderDetails) dto.ProviderDetails {
	return dto.ProviderDetails{
		AccountCount: model.AccountCount,
		ClusterCount: model.ClusterCount,
	}
}

func toScannerDTO(model inventory.Scanner) dto.Scanner {
	return dto.Scanner{
		LastScanTimestamp: model.LastScanTimestamp,
	}
}
