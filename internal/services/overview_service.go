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

	// TODO: Get Scanner data.
	// For now, returning empty data.

	return overview, nil
}

//nolint:cyclop
func (s *overviewServiceImpl) getProvidersSummary(ctx context.Context) (inventory.ProvidersSummary, error) {
	summary := inventory.ProvidersSummary{}
	const maxAccounts = 1000 // Assuming we have less than 1000 accounts
	opts := models.ListOptions{PageSize: maxAccounts, Offset: 0}

	accounts, _, err := s.accountRepo.ListAccounts(ctx, opts)
	if err != nil {
		return summary, err
	}

	clusters, _, err := s.clusterRepo.ListClusters(ctx, opts)
	if err != nil {
		return summary, err
	}

	for _, acc := range accounts {
		switch acc.Provider {
		case inventory.AWSProvider:
			summary.AWS.AccountCount++
		case inventory.GCPProvider:
			summary.GCP.AccountCount++
		case inventory.AzureProvider:
			summary.Azure.AccountCount++
		}
	}

	for _, cls := range clusters {
		switch cls.Provider {
		case inventory.AWSProvider:
			summary.AWS.ClusterCount++
		case inventory.GCPProvider:
			summary.GCP.ClusterCount++
		case inventory.AzureProvider:
			summary.Azure.ClusterCount++
		}
	}

	return summary, nil
}
