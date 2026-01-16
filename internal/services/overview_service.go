package services

import (
	"context"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// OverviewService defines the interface for overview-related business logic.
type OverviewService interface {
	GetOverview(ctx context.Context) (inventory.OverviewSummary, error)
}

var _ OverviewService = (*overviewServiceImpl)(nil)

type overviewServiceImpl struct {
	clusterRepo  repositories.ClusterRepository
	instanceRepo repositories.InstanceRepository
	accountRepo  repositories.AccountRepository
	// TODO: Add scanner repository when available
}

// NewOverviewService creates a new instance of OverviewService.
func NewOverviewService(clusterRepo repositories.ClusterRepository, instanceRepo repositories.InstanceRepository, accountRepo repositories.AccountRepository) OverviewService {
	return &overviewServiceImpl{
		clusterRepo:  clusterRepo,
		instanceRepo: instanceRepo,
		accountRepo:  accountRepo,
	}
}

// GetOverview retrieves all components of the inventory overview.
func (s *overviewServiceImpl) GetOverview(ctx context.Context) (inventory.OverviewSummary, error) {
	var overview inventory.OverviewSummary

	clusters, err := s.clusterRepo.GetClustersOverview(ctx)
	if err != nil {
		return inventory.OverviewSummary{}, fmt.Errorf("failed to get clusters overview: %w", err)
	}
	overview.Clusters = clusters

	instances, err := s.instanceRepo.GetInstancesOverview(ctx)
	if err != nil {
		return inventory.OverviewSummary{}, fmt.Errorf("failed to get instances overview: %w", err)
	}
	overview.Instances = instances

	providers, err := s.getProvidersSummary(ctx)
	if err != nil {
		return inventory.OverviewSummary{}, fmt.Errorf("failed to get providers summary: %w", err)
	}
	overview.Providers = providers

	scannerTimestamp, err := s.accountRepo.GetScannerTimestamp(ctx)
	if err != nil {
		return inventory.OverviewSummary{}, fmt.Errorf("failed to get scanner timestamp: %w", err)
	}
	overview.Scanner.LastScanTimestamp = scannerTimestamp

	return overview, nil
}

func (s *overviewServiceImpl) getProvidersSummary(ctx context.Context) (inventory.ProvidersSummary, error) {
	var err error
	var summary inventory.ProvidersSummary

	// Define a slice of providers with their names and corresponding summary buckets
	providers := []struct {
		name   string
		bucket *inventory.ProviderDetails
	}{
		{name: "AWS", bucket: &summary.AWS},
		{name: "GCP", bucket: &summary.GCP},
		{name: "Azure", bucket: &summary.Azure},
	}

	// Iterate over each provider to populate their summary details
	for _, p := range providers {
		// Set list options with filters for the current provider
		opts := models.ListOptions{
			PageSize: 0,
			Offset:   0,
			Filters: map[string]interface{}{
				"provider": p.name,
			},
		}

		// Count accounts for the current provider and handle errors
		if p.bucket.AccountCount, err = s.accountRepo.CountAccounts(ctx, opts); err != nil {
			return summary, err
		}

		// Count clusters for the current provider and handle errors
		if p.bucket.ClusterCount, err = s.clusterRepo.CountClusters(ctx, opts); err != nil {
			return summary, err
		}
	}

	// Return the populated summary of providers
	return summary, nil
}
