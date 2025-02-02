// Package main is the entry point for the ClusterIQ API server.
// It initializes the API server, sets up routes, and handles server lifecycle events.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/cmd/api/docs"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/middleware"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	// APITimeoutSeconds defines the default timeout in seconds for the API connection.
	// This value is used for graceful shutdowns and other timeout-related operations.
	APITimeoutSeconds = 60
)

var (
	// version reflects the current version of the API.
	// It is populated at build time using build flags.
	version string

	// commit reflects the git short-hash of the compiled version.
	// It provides traceability for the exact source code version used to build the binary.
	commit string
)

// APIServer represents the API server, including configuration, logger, router, and clients for gRPC and SQL.
type APIServer struct {
	cfg    *config.APIServerConfig // Configuration for the API server
	logger *zap.Logger             // Logger instance
	router *gin.Engine             // Gin router for handling HTTP requests
	server *http.Server            // HTTP server instance
	grpc   *APIGRPCClient          // gRPC client for communication with external services
	sql    *APISQLClient           // SQL client for database operations
}

// NewAPIServer initializes a new instance of the APIServer.
// It configures the Gin router, HTTP server, gRPC client, and SQL client.
//
// Parameters:
// - cfg: Configuration object for the API server.
// - logger: Logger instance for logging.
//
// Returns:
// - Pointer to the newly created APIServer.
func NewAPIServer(cfg *config.APIServerConfig, logger *zap.Logger) (*APIServer, error) {
	// Configuring GIN router
	router := setupGin(logger)
	// Configure HTTP server
	server := &http.Server{
		Addr:    cfg.ListenURL,
		Handler: router,
	}

	// Creating gRPC client
	gRPCClient, err := NewAPIGRPCClient(cfg.AgentURL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	// Creating DB client
	sqlClient, err := NewAPISQLClient(cfg.DBURL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQL client: %w", err)
	}

	return &APIServer{
		cfg:    cfg,
		logger: logger,
		server: server,
		router: router,
		grpc:   gRPCClient,
		sql:    sqlClient,
	}, nil
}

func setupGin(logger *zap.Logger) *gin.Engine {
	// TODO. Configure via env vars
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.SetCommonHeaders())
	// Configure Gin to use Zap
	router.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/api/v1/healthcheck"},
	}))
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	return router
}

func (a *APIServer) setupSwagger() {
	docs.SwaggerInfo.Title = "Cluster IP API doc"
	docs.SwaggerInfo.Description = "This the API of the ClusterIQ project"
	docs.SwaggerInfo.Version = "0.3"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

// Start starts the HTTP server in a goroutine
func (a *APIServer) Start() error {
	// Configure route
	a.setupRouter()
	// Configure Swagger
	a.setupSwagger()
	// Configure default middleware
	a.router.Use()
	a.logger.Info("==================== Starting ClusterIQ API ====================",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("api_url", a.cfg.ListenURL),
		zap.String("db_url", a.cfg.DBURL),
		zap.String("agent_url", a.cfg.AgentURL))

	// Start API
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal("Server listen and serve error", zap.Error(err))
			os.Exit(1)
		}
	}()
	a.logger.Info("API Ready to serve")
	return nil
}

// Run starts the server and handles graceful shutdown
func (a *APIServer) Run() error {
	if err := a.Start(); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	// Re-register to catch force shutdown
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit // Second signal
		os.Exit(1)
	}()

	return a.signalHandler(s)
}

// signalHandler handles OS signals for graceful server shutdown.  It shuts
// down the server when a SIGTERM signal is received. This function was included
// for better integration on K8s/OCP
//
// Parameters:
// - signal: The OS signal to handle.
func (a APIServer) signalHandler(signal os.Signal) error {
	if signal == syscall.SIGTERM {
		a.logger.Warn("SIGTERM signal received. Stopping ClusterIQ API server")
	} else {
		a.logger.Warn("Shutting down server...", zap.String("signal", signal.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), APITimeoutSeconds*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Fatal("API Shutdown error", zap.Error(err))
		return err
	}

	a.logger.Info("API server stopped")
	return nil
}

//	@title			ClusterIQ API
//	@version		0.3
//	@description	This is the API of the ClusterIQ cloud inventory software
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	ClusterIQ Team
//	@contact.email	cloud-native-team@redhat.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// Initialize logging configuration
	logger := ciqLogger.NewLogger()

	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	// Loading APIServer config
	cfg, err := config.LoadAPIServerConfig()
	if err != nil {
		logger.Error("Error loading APIServer config", zap.Error(err))
		return
	}

	// Initializing APIServer instance
	api, err := NewAPIServer(cfg, logger)
	if err != nil {
		logger.Fatal("Error creating API server", zap.Error(err))
		return
	}

	if err := api.Run(); err != nil {
		logger.Fatal("Error running API server", zap.Error(err))
		return
	}
}
