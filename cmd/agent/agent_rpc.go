package main

import (
	"context"
	"fmt"

	pb "github.com/RHEcosystemAppEng/cluster-iq/generated/agent"
	cexec "github.com/RHEcosystemAppEng/cluster-iq/internal/cloud_executors"
	"go.uber.org/zap"
)

const (
	// PowerOffClusterSuccessfully defines the success message format for powering off a cluster.
	PowerOffClusterSuccessfully = "Power Off for Cluster: %s(Acc: %s; Instances: %d) Successfull"
	// PowerOffClusterError defines the error message format for powering off a cluster.
	PowerOffClusterError = "Power Off for Cluster: %s(Acc: %s; Instances: %d) Failed"
	// PowerOnClusterSuccessfully defines the success message format for powering on a cluster.
	PowerOnClusterSuccessfully = "Power On for Cluster: %s(Acc: %s; Instances: %d) Successfull"
	// PowerOnClusterError defines the error message format for powering on a cluster.
	PowerOnClusterError = "Power On for Cluster: %s(Acc: %s; Instances: %d) Failed"
)

// GetExecutor retrieves the CloudExecutor associated with a given account name.
//
// Parameters:
// - accountName: The name of the account for which the executor is requested.
//
// Returns:
// - cexec.CloudExecutor: The executor for the specified account.
// - error: An error if no executor is found for the given account.
func (a *AgentService) GetExecutor(accountName string) (cexec.CloudExecutor, error) {
	exec, ok := a.executors[accountName]
	if !ok {
		return nil, fmt.Errorf("There's no Executor available for the requested account")
	}
	return exec, nil
}

// PowerOnCluster handles a gRPC request to power on a cluster.
//
// Parameters:
// - ctx: The context for the gRPC request.
// - req: The request object containing the cluster ID, account name, region, and instances.
//
// Returns:
// - *pb.PowerOnClusterResponse: The response object containing a success or error message.
// - error: An error if the operation fails.
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
		zap.Int("instances_num", len(req.InstancesIdList)),
	)

	exec.PowerOnCluster(req.InstancesIdList)
	return &pb.PowerOnClusterResponse{
		Error:   0,
		Message: fmt.Sprintf(PowerOnClusterSuccessfully, req.ClusterId, req.AccountName, len(req.InstancesIdList)),
	}, nil
}

// PowerOffCluster handles a gRPC request to power Off a cluster.
//
// Parameters:
// - ctx: The context for the gRPC request.
// - req: The request object containing the cluster ID, account name, region, and instances.
//
// Returns:
// - *pb.PowerOffClusterResponse: The response object containing a success or error message.
// - error: An error if the operation fails.
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
