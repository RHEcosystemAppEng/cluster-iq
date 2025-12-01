package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/RHEcosystemAppEng/cluster-iq/generated/agent"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// InstantAgentService represents the main structure for managing cloud executors and configuration.
// It also embeds the gRPC server interface for handling gRPC requests.
type InstantAgentService struct {
	// Basic properties
	cfg *config.InstantAgentServiceConfig
	AgentService
	// gRPC properties
	pb.UnimplementedAgentServiceServer
	grpcServer *grpc.Server
	listener   net.Listener
}

// NewInstantAgentService creates and initializes a new AgentService instance for serving gRPC requests
//
// Parameters:
//   - cfg: Pointer to AgentServiceConfig containing the configuration details.
//   - logger: Pointer to zap.Logger for logging.
//
// Returns:
//   - *InstantAgentService: A pointer to the newly created AgentService instance.
func NewInstantAgentService(cfg *config.InstantAgentServiceConfig, actionsChannel chan<- actions.Action, wg *sync.WaitGroup, logger *zap.Logger) *InstantAgentService {
	// Listener config
	lis, err := net.Listen("tcp", cfg.ListenURL)
	if err != nil {
		logger.Error("Error initializing gRPC AgentService on ClusterIQ Agent", zap.Error(err))
		return nil
	}

	// Initializing gRPC server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))
	reflection.Register(grpcServer)

	ias := &InstantAgentService{
		cfg: cfg,
		AgentService: AgentService{
			logger:         logger,
			wg:             wg,
			actionsChannel: actionsChannel,
		},
		grpcServer: grpcServer,
		listener:   lis,
	}

	// Registering Agent service on gRPC server
	pb.RegisterAgentServiceServer(ias.grpcServer, ias)

	ias.logger.Info("gRPC ClusterIQ Agent initialization successfully",
		zap.String("listen_url", ias.cfg.ListenURL),
		zap.String("version", version),
		zap.String("commit", commit))

	return ias
}

func (i *InstantAgentService) Start() error {
	defer i.wg.Done()

	i.logger.Info("Starting AgentService")

	// Serving gRPC
	if err := i.grpcServer.Serve(i.listener); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
		return err
	}
	return nil
}

// PowerOnCluster handles a gRPC request to power on a cluster.
//
// Parameters:
// - ctx: The context for the gRPC request. (Not used)
// - req: The request object containing the cluster ID, account name, region, and instances.
//
// Returns:
// - *pb.PowerOnClusterResponse: The response object containing a success or error message.
// - error: An error if the operation fails.
func (i *InstantAgentService) PowerOnCluster(_ context.Context, req *pb.PowerOnClusterRequest) (*pb.PowerOnClusterResponse, error) {
	i.logger.Warn("Powering On Cluster",
		zap.String("account_name", req.AccountId),
		zap.String("region", req.Region),
		zap.String("cluster_id", req.ClusterId),
		zap.Strings("instances", req.InstancesIdList),
		zap.Int("instances_num", len(req.InstancesIdList)),
	)

	action := actions.NewInstantAction(
		actions.PowerOnCluster,
		*actions.NewActionTarget(
			req.AccountId,
			req.Region,
			req.ClusterId,
			req.InstancesIdList,
		),
		"Pending",
		true,
	)

	i.actionsChannel <- action

	return &pb.PowerOnClusterResponse{
		Error:   0,
		Message: fmt.Sprintf(PowerOnClusterSuccessfully, req.ClusterId, req.AccountId, len(req.InstancesIdList)),
	}, nil
}

// PowerOffCluster handles a gRPC request to power Off a cluster.
//
// Parameters:
// - ctx: The context for the gRPC request. (Not used)
// - req: The request object containing the cluster ID, account name, region, and instances.
//
// Returns:
// - *pb.PowerOffClusterResponse: The response object containing a success or error message.
// - error: An error if the operation fails.
func (i *InstantAgentService) PowerOffCluster(_ context.Context, req *pb.PowerOffClusterRequest) (*pb.PowerOffClusterResponse, error) {
	i.logger.Warn("Powering Off Cluster",
		zap.String("account_name", req.AccountId),
		zap.String("region", req.Region),
		zap.String("cluster_id", req.ClusterId),
		zap.Strings("instances", req.InstancesIdList),
		zap.Int("instances_num", len(req.InstancesIdList)))

	action := actions.NewInstantAction(
		actions.PowerOffCluster,
		*actions.NewActionTarget(
			req.AccountId,
			req.Region,
			req.ClusterId,
			req.InstancesIdList,
		),
		"Pending",
		true,
	)

	i.actionsChannel <- action

	return &pb.PowerOffClusterResponse{
		Error:   0,
		Message: fmt.Sprintf(PowerOffClusterSuccessfully, req.ClusterId, req.AccountId, len(req.InstancesIdList)),
	}, nil
}
