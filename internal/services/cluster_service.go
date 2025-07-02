package services

import (
	"context"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/clients"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// ClusterService defines the interface for cluster-related business logic.
type ClusterService interface {
	List(ctx context.Context, options repositories.ListOptions) ([]inventory.Cluster, int, error)
	Get(ctx context.Context, id string) (inventory.Cluster, error)
	GetSummary(ctx context.Context) (inventory.ClustersSummary, error)
	PowerOn(ctx context.Context, clusterID string) error
	PowerOff(ctx context.Context, clusterID string) error
	Create(ctx context.Context, cluster inventory.Cluster) error
	Delete(ctx context.Context, clusterID string) error
	GetTags(ctx context.Context, clusterID string) ([]inventory.Tag, error)
	Update(ctx context.Context, cluster inventory.Cluster) error
}

var _ ClusterService = (*clusterServiceImpl)(nil)

type clusterServiceImpl struct {
	repo        repositories.ClusterRepository
	agentClient clients.AgentClient
}

// NewClusterService creates a new instance of ClusterService.
func NewClusterService(repo repositories.ClusterRepository, agentClient clients.AgentClient) ClusterService {
	return &clusterServiceImpl{
		repo:        repo,
		agentClient: agentClient,
	}
}

// List retrieves a paginated list of clusters based on the provided options.
func (s *clusterServiceImpl) List(ctx context.Context, options repositories.ListOptions) ([]inventory.Cluster, int, error) {
	return s.repo.ListClusters(ctx, options)
}

// Get retrieves a single cluster by its ID.
func (s *clusterServiceImpl) Get(ctx context.Context, id string) (inventory.Cluster, error) {
	return s.repo.GetClusterByID(ctx, id)
}

// GetSummary retrieves a summary of cluster counts by status.
func (s *clusterServiceImpl) GetSummary(ctx context.Context) (inventory.ClustersSummary, error) {
	return s.repo.GetClustersOverview(ctx)
}

// PowerOn sends a request to power on a cluster.
func (s *clusterServiceImpl) PowerOn(ctx context.Context, clusterID string) error {
	cluster, err := s.repo.GetClusterByID(ctx, clusterID)
	if err != nil {
		return err
	}
	instances, err := s.repo.GetInstancesOnCluster(ctx, clusterID)
	if err != nil {
		return err
	}
	instanceIDs := make([]string, len(instances))
	for i, inst := range instances {
		instanceIDs[i] = inst.ID
	}

	req := &clients.ClusterStatusChangeRequest{
		AccountName:     cluster.AccountName,
		Region:          cluster.Region,
		ClusterID:       cluster.ID,
		InstancesIdList: instanceIDs,
	}

	return s.agentClient.PowerOnCluster(ctx, req)
}

// PowerOff sends a request to power off a cluster.
func (s *clusterServiceImpl) PowerOff(ctx context.Context, clusterID string) error {
	cluster, err := s.repo.GetClusterByID(ctx, clusterID)
	if err != nil {
		return err
	}
	instances, err := s.repo.GetInstancesOnCluster(ctx, clusterID)
	if err != nil {
		return err
	}
	instanceIDs := make([]string, len(instances))
	for i, inst := range instances {
		instanceIDs[i] = inst.ID
	}

	req := &clients.ClusterStatusChangeRequest{
		AccountName:     cluster.AccountName,
		Region:          cluster.Region,
		ClusterID:       cluster.ID,
		InstancesIdList: instanceIDs,
	}

	return s.agentClient.PowerOffCluster(ctx, req)
}

// Create creates a new cluster.
func (s *clusterServiceImpl) Create(ctx context.Context, cluster inventory.Cluster) error {
	return s.repo.CreateCluster(ctx, cluster)
}

// Delete deletes a cluster by its ID.
func (s *clusterServiceImpl) Delete(ctx context.Context, clusterID string) error {
	return s.repo.DeleteCluster(ctx, clusterID)
}

// GetTags retrieves all tags for a specific cluster.
func (s *clusterServiceImpl) GetTags(ctx context.Context, clusterID string) ([]inventory.Tag, error) {
	return s.repo.GetClusterTags(ctx, clusterID)
}

// Update updates an existing cluster.
func (s *clusterServiceImpl) Update(ctx context.Context, cluster inventory.Cluster) error {
	return s.repo.UpdateCluster(ctx, cluster)
}
