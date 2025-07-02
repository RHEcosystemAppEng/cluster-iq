package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/clients"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

type ClusterService interface {
	PowerOn(ctx context.Context, clusterID string) error
}

type clusterServiceImpl struct {
	repo        repositories.ClusterRepository
	agentClient clients.AgentClient
}

func NewClusterService(repo repositories.ClusterRepository, agentClient clients.AgentClient) ClusterService {
	return &clusterServiceImpl{repo: repo, agentClient: agentClient}
}

func (s *clusterServiceImpl) PowerOn(ctx context.Context, clusterID string) error {
	cluster, err := s.repo.GetClusterByID(clusterID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return fmt.Errorf("cluster with ID '%s' not found", clusterID)
		}
		return fmt.Errorf("repository failed to get cluster: %w", err)
	}

	instances, err := s.repo.GetInstancesOnCluster(clusterID)
	if err != nil {
		return fmt.Errorf("repository failed to get instances: %w", err)
	}

	instanceIDs := make([]string, 0, len(instances))
	for _, inst := range instances {
		instanceIDs = append(instanceIDs, inst.ID)
	}

	agentRequest := &clients.ClusterStatusChangeRequest{
		AccountName:     cluster.AccountName,
		Region:          cluster.Region,
		ClusterID:       cluster.ID,
		InstancesIdList: instanceIDs,
	}

	if err := s.agentClient.PowerOnCluster(ctx, agentRequest); err != nil {
		return fmt.Errorf("agent client failed to power on cluster: %w", err)
	}
	return nil
}
