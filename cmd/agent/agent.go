// Package main implements the ClusterIQ Agent application.
//
// The ClusterIQ Agent is responsible for managing cloud provider accounts,
// initializing cloud executors, and exposing gRPC services for external interaction.
// It serves as a key component of the ClusterIQ system, enabling inventory management
// and operations on cloud resources.
//
// Features of the ClusterIQ Agent:
// - Manages configurations for multiple cloud providers.
// - Initializes and maintains executors for cloud provider accounts.
// - Provides a gRPC service interface for interacting with the ClusterIQ system.
// - Logs detailed information about operations for debugging and monitoring.
//
// The application uses gRPC as the communication protocol and supports extensible
// cloud executor implementations for AWS, GCP, and Azure.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"

	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	// logger is a shared logging instance used across the entire Agent application.
	logger *zap.Logger
	// version reflects the current version of the Agent.
	// It is populated at build time using build flags.
	version string
	// commit reflects the git short-hash of the compiled version.
	// It provides traceability for the exact source code version used to build the binary.
	commit string
)

// init initializes the global logger configuration.
// This is automatically invoked before the main function.
func init() {
	// Initialize logging configuration
	logger = ciqLogger.NewLogger()
}

// Agent represents the main structure for the ClusterIQ Agent. It includes the
// AgentService for implementing the gRPC server and the AgentCron for getting
// and executing the Scheduled actions to be run by the AgentService when it's
// needed
type Agent struct {
	cfg            *config.AgentConfig
	ias            *InstantAgentService
	cas            *ScheduleAgentService
	eas            *ExecutorAgentService
	actionsChannel chan actions.Action
	logger         *zap.Logger
	wg             *sync.WaitGroup
}

func NewAgent(cfg *config.AgentConfig, logger *zap.Logger) *Agent {
	var ch chan actions.Action
	var wg sync.WaitGroup

	ch = make(chan actions.Action)

	// Creating InstantAgentService (gRPC)
	ias := NewInstantAgentService(&cfg.InstantAgentServiceConfig, ch, &wg, logger)
	if ias == nil {
		logger.Error("Cannot create InstantAgentService")
		return nil
	}

	// Creating ScheduleAgentService (scheduled actions)
	cas := NewScheduleAgentService(&cfg.ScheduleAgentServiceConfig, ch, &wg, logger)
	if ias == nil {
		logger.Error("Cannot create CronAgentService")
		return nil
	}

	// Creating ExecutorAgentService (executing actions)
	eas := NewExecutorAgentService(&cfg.ExecutorAgentServiceConfig, ch, &wg, logger)
	if eas == nil {
		logger.Error("Cannot create ExecutorAgentService")
		return nil
	}

	return &Agent{
		cfg:            cfg,
		ias:            ias,
		cas:            cas,
		eas:            eas,
		actionsChannel: ch,
		logger:         logger,
		wg:             &wg,
	}
}

// TODO DOC
func (a *Agent) StartAgentServices() error {
	var err error
	errChan := make(chan error, 3)

	// Starting InstantAgentService
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err = a.ias.Start(); err != nil {
			errChan <- fmt.Errorf("Instant AgentService (gRPC) failed: %w", err)
			return
		}
		a.logger.Info("Instant Agent Service (gRPC) finished")
	}()

	// Starting ScheduledAgentService
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err = a.cas.Start(); err != nil {
			errChan <- fmt.Errorf("Scheduled Agent Service failed: %w", err)
			return
		}
		a.logger.Info("Scheduled Agent Service finished")
	}()

	// Starting ExecutorAgentService
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err = a.eas.Start(); err != nil {
			errChan <- fmt.Errorf("Executor Agent Service failed: %w", err)
			return
		}
		a.logger.Info("Executor Agent Service finished")
	}()

	a.logger.Info("ClusterIQ Agent Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	// Re-register to catch force shutdown
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	if err := a.signalHandler(s); err != nil {
		a.logger.Warn("Signal Captured", zap.Error(err))
	}

	go func() {
		<-quit // Second signal
		os.Exit(1)
	}()

	// Wait for goroutines to finish
	a.wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

// LoggingInterceptor is a gRPC interceptor that logs information about incoming requests and their responses.
//
// It logs details such as the client's IP address, the invoked method, and any errors that occur during method execution.
// This interceptor can be used to enhance visibility and debugging in gRPC server operations.
//
// Parameters:
// - ctx: The context of the gRPC request.
// - req: The incoming request payload.
// - info: Metadata about the invoked gRPC method (e.g., method name).
// - handler: The actual handler function that processes the request.
//
// Returns:
// - An interface{} representing the response from the handler.
// - An error if the handler or any other operation fails.
func LoggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		logger.Info("Client connected", zap.String("ip", p.Addr.String()))
	}

	log.Printf("Invoked method: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error in method %s: %v", info.FullMethod, err)
	} else {
		log.Printf("Method %s executed successfully", info.FullMethod)
	}

	return resp, err
}

// signalHandler handles OS signals for graceful server shutdown.  It shuts
// down the server when a SIGTERM signal is received. This function was included
// for better integration on K8s/OCP
//
// Parameters:
// - signal: The OS signal to handle.
func (a Agent) signalHandler(signal os.Signal) error {
	if signal == syscall.SIGTERM {
		a.logger.Warn("SIGTERM signal received. Stopping ClusterIQ Agent server")
	} else {
		a.logger.Warn("Shutting down server...", zap.String("signal", signal.String()))
	}

	for _, item := range a.cas.schedule {
		cancel := item.cancel
		action := item.action
		a.logger.Warn("Cancelling ScheduledAction", zap.String("action_id", action.GetID()))
		cancel()
	}
	close(a.actionsChannel)

	a.logger.Info("ClusterIQ server stopped")
	os.Exit(0)
	return nil
}

func (a Agent) Start() error {
	a.logger.Info("==================== Starting ClusterIQ Agent ====================",
		zap.String("version", version),
		zap.String("commit", commit),
	)

	return a.StartAgentServices()
}

// main is the entry point for the ClusterIQ Agent application.
// It initializes the Agent, loads configuration, creates cloud executors, and starts the gRPC server.
func main() {
	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	var err error

	// Loading AgentService configuration
	cfg, err := config.LoadAgentConfig()
	if err != nil {
		logger.Fatal("Failed to load Agent config", zap.Error(err))
	}

	// Creating AgentService with the specified configuration
	agent := NewAgent(cfg, logger)
	if agent == nil {
		logger.Error("Error during AgentService setup. Aborting Agent")
		os.Exit(-1)
	}

	// Starting Agent
	err = agent.Start()

	agent.logger.Info("ClusterIQ Agent Finished")
	if err != nil {
		agent.logger.Error("Error starting Agent Services", zap.Error(err))
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
