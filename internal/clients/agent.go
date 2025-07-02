package clients

import (
	"context"
	pb "github.com/RHEcosystemAppEng/cluster-iq/generated/agent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TOdO, review timeout
const (
	// grpcTimeoutSeconds defines the timeout in seconds for gRPC operations.
	grpcTimeoutSeconds = 10
)

// ClusterStatusChangeRequest represents the request to the gRPC Agent for powering on/off clusters.
// It includes details such as the account name, region, cluster ID, and the list of instance IDs associated with the cluster.
type ClusterStatusChangeRequest struct {
	AccountName     string   // The name of the account associated with the cluster.
	Region          string   // The AWS region where the cluster is located.
	ClusterID       string   // The unique identifier of the cluster.
	InstancesIdList []string // A list of instance IDs belonging to the cluster.
}

type AgentClient interface {
	PowerOnCluster(ctx context.Context, request *ClusterStatusChangeRequest) error
	PowerOffCluster(ctx context.Context, request *ClusterStatusChangeRequest) error
}

// GRPCAgentClient manages the gRPC client connection and operations for the API server.
// It provides methods to interact with the Agent service via gRPC.
type GRPCAgentClient struct {
	// Client is the gRPC client used to communicate with the Agent service.
	Client pb.AgentServiceClient
	// logger is used for logging gRPC operations and errors.
	logger *zap.Logger
}

// NewGRPCAgentClient initializes and returns a new GRPCAgentClient.
// It establishes a connection to the Agent service and sets up the gRPC client.
//
// Parameters:
// - agentURL: The URL of the Agent service to connect to.
// - logger: Logger instance for logging.
//
// Returns:
// - A pointer to the initialized GRPCAgentClient.
// - An error if the connection cannot be established.
func NewGRPCAgentClient(agentURL string, logger *zap.Logger) (*GRPCAgentClient, error) {
	// Initializing gRPC Client
	conn, err := grpc.NewClient(agentURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCAgentClient{
		Client: pb.NewAgentServiceClient(conn),
		logger: logger,
	}, nil
}

// PowerOffCluster sends a gRPC request to power off a cluster by the given ClusterID.
// It logs the details of the request and the response received.
//
// Parameters:
// - request: A ClusterStatusChangeRequest containing details about the cluster to power off.
//
// Returns:
// - An error if the gRPC call fails or the request cannot be completed.
func (c *GRPCAgentClient) PowerOffCluster(ctx context.Context, request *ClusterStatusChangeRequest) error {
	// Creating PowerOffClusterRequest
	rpcRequest := &pb.PowerOffClusterRequest{
		AccountName:     request.AccountName,
		Region:          request.Region,
		ClusterId:       request.ClusterID,
		InstancesIdList: request.InstancesIdList,
	}

	// Logging the request details
	c.logger.Info("Powering off Cluster",
		zap.String("account_name", rpcRequest.AccountName),
		zap.String("cluster_id", rpcRequest.ClusterId),
		zap.String("region", rpcRequest.Region),
		zap.Strings("instances", rpcRequest.InstancesIdList),
		zap.Int("instances_count", len(rpcRequest.InstancesIdList)),
	)

	// Sending the PowerOffCluster request
	resp, err := c.Client.PowerOffCluster(ctx, rpcRequest)
	if err != nil {
		return err
	}
	c.logger.Info("Response from PowerOffCluster", zap.String("response", resp.Message))
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
// TODO
func (c *GRPCAgentClient) PowerOnCluster(ctx context.Context, request *ClusterStatusChangeRequest) error {
	// Creating PowerOnClusterRequest
	rpcRequest := &pb.PowerOnClusterRequest{
		AccountName:     request.AccountName,
		Region:          request.Region,
		ClusterId:       request.ClusterID,
		InstancesIdList: request.InstancesIdList,
	}

	// Logging the request details
	c.logger.Info("Powering On Cluster",
		zap.String("account_name", rpcRequest.AccountName),
		zap.String("cluster_id", rpcRequest.ClusterId),
		zap.String("region", rpcRequest.Region),
		zap.Strings("instances", rpcRequest.InstancesIdList),
		zap.Int("instances_count", len(rpcRequest.InstancesIdList)),
	)

	// Sending the PowerOnCluster request
	resp, err := c.Client.PowerOnCluster(ctx, rpcRequest)
	if err != nil {
		return err
	}
	c.logger.Info("Response from PowerOnCluster", zap.String("response", resp.Message))
	return nil
}
