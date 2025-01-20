package main

import (
	"context"
	"fmt"

	pb "github.com/RHEcosystemAppEng/cluster-iq/generated/agent"
	cexec "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_executors"
	"go.uber.org/zap"
)

const (
	PowerOffClusterSuccessfully = "Power Off for Cluster: %s(Acc: %s; Instances: %d) Successfull"
	PowerOffClusterError        = "Power Off for Cluster: %s(Acc: %s; Instances: %d) Failed"
	PowerOnClusterSuccessfully  = "Power On for Cluster: %s(Acc: %s; Instances: %d) Successfull"
	PowerOnClusterError         = "Power On for Cluster: %s(Acc: %s; Instances: %d) Failed"
)

func (a *AgentService) GetExecutor(accountName string) (cexec.CloudExecutor, error) {
	exec, ok := a.executors[accountName]
	if !ok {
		return nil, fmt.Errorf("There's no Executor available for the requested account")
	}
	return exec, nil
}

// PowerOnCluster gRPC function for powering off a cluster
func (a *AgentService) PowerOnCluster(ctx context.Context, req *pb.PowerOnClusterRequest) (*pb.PowerOnClusterResponse, error) {
	a.logger.Debug("Received PowerOnCluster Request", zap.String("cluster_id", req.ClusterId), zap.String("accound_name", req.AccountName), zap.Int("instances", len(req.InstancesIdList)))

	// Getting Executor for the requested account
	exec, err := a.GetExecutor(req.AccountName)
	if err != nil {
		a.logger.Error("Can't get executor for PowerOnCluster", zap.String("account_name", req.AccountName), zap.Error(err))
		return &pb.PowerOnClusterResponse{
			Error:   1,
			Message: fmt.Sprintf(PowerOnClusterError, req.ClusterId, req.AccountName, len(req.InstancesIdList)),
		}, err
	}

	// Configuring Region
	if err := exec.SetRegion(req.Region); err != nil {
		a.logger.Error("Error configuring region before executing a PowerOnCluster", zap.Error(err))
		return &pb.PowerOnClusterResponse{
			Error:   1,
			Message: fmt.Sprintf(PowerOnClusterError, req.ClusterId, req.AccountName, len(req.InstancesIdList)),
		}, err
	}

	// PowerOn
	logger.Warn("Powering On Cluster",
		zap.String("accound_id", req.AccountName),
		zap.String("region", req.Region),
		zap.String("cluster_id", req.ClusterId),
		zap.Strings("instances", req.InstancesIdList),
		zap.Int("instances_num", len(req.InstancesIdList)))
	exec.PowerOnCluster(req.InstancesIdList)
	return &pb.PowerOnClusterResponse{
		Error:   0,
		Message: fmt.Sprintf(PowerOnClusterSuccessfully, req.ClusterId, req.AccountName, len(req.InstancesIdList)),
	}, nil
}

// PowerOffCluster gRPC function for powering on a cluster
func (a *AgentService) PowerOffCluster(ctx context.Context, req *pb.PowerOffClusterRequest) (*pb.PowerOffClusterResponse, error) {
	a.logger.Debug("Received PowerOffCluster Request", zap.String("cluster_id", req.ClusterId), zap.String("accound_name", req.AccountName), zap.Int("instances", len(req.InstancesIdList)))

	// Getting Executor for the requested account
	exec, err := a.GetExecutor(req.AccountName)
	if err != nil {
		a.logger.Error("Can't get executor for PowerOffCluster", zap.String("account_name", req.AccountName), zap.Error(err))
		return &pb.PowerOffClusterResponse{
			Error:   1,
			Message: fmt.Sprintf(PowerOffClusterError, req.ClusterId, req.AccountName, len(req.InstancesIdList)),
		}, err
	}

	// Configuring Region
	if err := exec.SetRegion(req.Region); err != nil {
		a.logger.Error("Error configuring region before executing a PowerOffCluster", zap.Error(err))
		return &pb.PowerOffClusterResponse{
			Error:   1,
			Message: fmt.Sprintf(PowerOffClusterError, req.ClusterId, req.AccountName, len(req.InstancesIdList)),
		}, err
	}

	// PowerOff
	logger.Warn("Powering On Cluster",
		zap.String("accound_id", req.AccountName),
		zap.String("region", req.Region),
		zap.String("cluster_id", req.ClusterId),
		zap.Strings("instances", req.InstancesIdList),
		zap.Int("instances_num", len(req.InstancesIdList)))
	exec.PowerOffCluster(req.InstancesIdList)
	return &pb.PowerOffClusterResponse{
		Error:   0,
		Message: fmt.Sprintf(PowerOffClusterSuccessfully, req.ClusterId, req.AccountName, len(req.InstancesIdList)),
	}, nil
}
