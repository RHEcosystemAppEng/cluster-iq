package services

import (
	"context"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/clients"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// ClusterService defines the interface for cluster-related business logic.
type ClusterService interface {
	List(ctx context.Context, options models.ListOptions) ([]db.ClusterDBResponse, int, error)
	Get(ctx context.Context, clusterID string) (*db.ClusterDBResponse, error)
	GetInstances(ctx context.Context, clusterID string) ([]db.InstanceDBResponse, error)
	GetSummary(ctx context.Context) (inventory.ClustersSummary, error)
	PowerOn(ctx context.Context, clusterID string) error
	PowerOff(ctx context.Context, clusterID string) error
	Create(ctx context.Context, clusters []inventory.Cluster) error
	Delete(ctx context.Context, clusterID string) error
	GetTags(ctx context.Context, clusterID string) ([]db.TagDBResponse, error)
	Update(ctx context.Context, cluster dto.ClusterDTORequest) error
}

var _ ClusterService = (*clusterServiceImpl)(nil)

type clusterServiceImpl struct {
	repo                repositories.ClusterRepository
	agentClient         clients.AgentClient
	agentRequestTimeout time.Duration
}

// ClusterServiceOptions contains configurable parameters for the ClusterService.
type ClusterServiceOptions struct {
	AgentRequestTimeout int
}

// NewClusterService creates a new instance of ClusterService.
func NewClusterService(repo repositories.ClusterRepository, agentClient clients.AgentClient, opts ClusterServiceOptions) ClusterService {
	return &clusterServiceImpl{
		repo:                repo,
		agentClient:         agentClient,
		agentRequestTimeout: time.Duration(opts.AgentRequestTimeout) * time.Second,
	}
}

// List retrieves a paginated list of clusters based on the provided options.
func (s *clusterServiceImpl) List(ctx context.Context, options models.ListOptions) ([]db.ClusterDBResponse, int, error) {
	return s.repo.ListClusters(ctx, options)
}

// Get retrieves a single cluster by its ID.
func (s *clusterServiceImpl) Get(ctx context.Context, clusterID string) (*db.ClusterDBResponse, error) {
	return s.repo.GetClusterByID(ctx, clusterID)
}

// Get retrieves a single cluster by its ID.
func (s *clusterServiceImpl) GetInstances(ctx context.Context, clusterID string) ([]db.InstanceDBResponse, error) {
	return s.repo.GetInstancesOnCluster(ctx, clusterID)
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
		instanceIDs[i] = inst.InstanceID
	}

	req := &clients.ClusterStatusChangeRequest{
		AccountID:       cluster.AccountID,
		Region:          cluster.Region,
		ClusterID:       cluster.ClusterID,
		InstancesIdList: instanceIDs,
	}

	ctx, cancel := context.WithTimeout(ctx, s.agentRequestTimeout)
	defer cancel()

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
		instanceIDs[i] = inst.InstanceID
	}

	req := &clients.ClusterStatusChangeRequest{
		AccountID:       cluster.AccountID,
		Region:          cluster.Region,
		ClusterID:       cluster.ClusterID,
		InstancesIdList: instanceIDs,
	}

	ctx, cancel := context.WithTimeout(ctx, s.agentRequestTimeout)
	defer cancel()

	return s.agentClient.PowerOffCluster(ctx, req)
}

// Create creates new clusters.
func (s *clusterServiceImpl) Create(ctx context.Context, clusters []inventory.Cluster) error {
	return s.repo.CreateClusters(ctx, clusters)
}

// Delete deletes a cluster by its ID.
func (s *clusterServiceImpl) Delete(ctx context.Context, clusterID string) error {
	return s.repo.DeleteCluster(ctx, clusterID)
}

// GetTags retrieves all tags for a specific cluster.
func (s *clusterServiceImpl) GetTags(ctx context.Context, clusterID string) ([]db.TagDBResponse, error) {
	return s.repo.GetClusterTags(ctx, clusterID)
}

// Update updates an existing cluster.
func (s *clusterServiceImpl) Update(ctx context.Context, cluster dto.ClusterDTORequest) error {
	return s.repo.UpdateCluster(ctx, cluster)
}
