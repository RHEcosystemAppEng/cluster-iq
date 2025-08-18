package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/RHEcosystemAppEng/cluster-iq/generated/agent"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// grpcTimeoutSeconds defines the timeout in seconds for gRPC operations.
	grpcTimeoutSeconds = 10
)

var (
	DefaultInstantActionDescription = "InstantAction"
)

// APIGRPCClient manages the gRPC client connection and operations for the API server.
// It provides methods to interact with the Agent service via gRPC.
type APIGRPCClient struct {
	// Client is the gRPC client used to communicate with the Agent service.
	Client pb.AgentServiceClient
	// CTX is the context used for gRPC operations.
	CTX context.Context
	// Cancel is the function to cancel the gRPC context.
	Cancel context.CancelFunc
	// logger is used for logging gRPC operations and errors.
	logger *zap.Logger
}

// NewAPIGRPCClient initializes and returns a new APIGRPCClient.
// It establishes a connection to the Agent service and sets up the gRPC client.
//
// Parameters:
// - agentURL: The URL of the Agent service to connect to.
// - logger: Logger instance for logging.
//
// Returns:
// - A pointer to the initialized APIGRPCClient.
// - An error if the connection cannot be established.
func NewAPIGRPCClient(agentURL string, logger *zap.Logger) (*APIGRPCClient, error) {
	// Initializing gRPC Client
	conn, err := grpc.NewClient(agentURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeoutSeconds*time.Second)

	return &APIGRPCClient{
		Client: pb.NewAgentServiceClient(conn),
		CTX:    ctx,
		Cancel: cancel,
		logger: logger,
	}, nil
}

func (a APIGRPCClient) ProcessInstantAction(action *actions.InstantAction) error {
	if action.GetDescription() == nil {
		action.Description = &DefaultInstantActionDescription
	}
	switch action.Operation {
	case actions.PowerOffCluster:
		return a.PowerOffCluster(action)
	case actions.PowerOnCluster:
		return a.PowerOnCluster(action)
	default:
		return fmt.Errorf("received InstantAction with unknown Operation")
	}
}

// PowerOffCluster sends a gRPC request to power off a cluster by the given ClusterID.
// It logs the details of the request and the response received.
//
// Parameters:
// - request: A ClusterStatusChangeRequest containing details about the cluster to power off.
//
// Returns:
// - An error if the gRPC call fails or the request cannot be completed.
func (a APIGRPCClient) PowerOffCluster(action *actions.InstantAction) error {
	// Creating PowerOffClusterRequest
	rpcRequest := &pb.PowerOffClusterRequest{
		AccountName:     action.GetTarget().AccountName,
		Region:          action.GetTarget().Region,
		ClusterId:       action.GetTarget().ClusterID,
		InstancesIdList: action.GetTarget().Instances,
		Requester:       action.GetRequester(),
		Description:     *action.GetDescription(),
	}

	// Logging the request details
	a.logger.Info("Powering off Cluster",
		zap.String("account_name", rpcRequest.AccountName),
		zap.String("cluster_id", rpcRequest.ClusterId),
		zap.String("region", rpcRequest.Region),
		zap.Strings("instances", rpcRequest.InstancesIdList),
		zap.Int("instances_count", len(rpcRequest.InstancesIdList)),
	)

	// Sending the PowerOffCluster request
	resp, err := a.Client.PowerOffCluster(context.Background(), rpcRequest)
	if err != nil {
		return err
	}
	a.logger.Info("Response from PowerOffCluster", zap.String("response", resp.Message))
	return nil
}

// PowerOnCluster sends a gRPC request to power on a cluster by the given ClusterID.
// It logs the details of the request and the response received.
//
// Parameters:
// - request: A ClusterStatusChangeRequest containing details about the cluster to power on.
//
// Returns:
// - An error if the gRPC call fails or the request cannot be completed.
func (a APIGRPCClient) PowerOnCluster(action *actions.InstantAction) error {
	// Creating PowerOnClusterRequest
	rpcRequest := &pb.PowerOnClusterRequest{
		AccountName:     action.GetTarget().AccountName,
		Region:          action.GetTarget().Region,
		ClusterId:       action.GetTarget().ClusterID,
		InstancesIdList: action.GetTarget().Instances,
		Requester:       action.GetRequester(),
		Description:     *action.GetDescription(),
	}

	// Logging the request details
	a.logger.Info("Powering On Cluster",
		zap.String("account_name", rpcRequest.AccountName),
		zap.String("cluster_id", rpcRequest.ClusterId),
		zap.String("region", rpcRequest.Region),
		zap.Strings("instances", rpcRequest.InstancesIdList),
		zap.Int("instances_count", len(rpcRequest.InstancesIdList)),
	)

	// Sending the PowerOnCluster request
	resp, err := a.Client.PowerOnCluster(context.Background(), rpcRequest)
	if err != nil {
		return err
	}
	a.logger.Info("Response from PowerOnCluster", zap.String("response", resp.Message))
	return nil
}
