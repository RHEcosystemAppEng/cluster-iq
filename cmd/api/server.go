package main

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/handlers"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/clients"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/middleware"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
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
	logger *zap.Logger  // Logger instance
	server *http.Server // HTTP server instance
	// eventService *events.EventService    // Service for handling audit logs
}

// NewAPIServer initializes a new instance of the APIServer.
// It configures the Gin router and HTTP server.
//
// Parameters:
// - listenAddr: Listen address.
// - logger: Logger instance for logging.
//
// Returns:
// - Pointer to the newly created APIServer.
func NewAPIServer(listenAddr string, logger *zap.Logger, engine *gin.Engine) *APIServer {
	return &APIServer{
		logger: logger,
		server: &http.Server{
			Addr:    listenAddr,
			Handler: engine,
		},
	}
}

func setupGin(logger *zap.Logger) *gin.Engine {
	// TODO. Configure via env vars
	gin.SetMode(gin.ReleaseMode)
	rtr := gin.New()
	// Configure default middleware
	rtr.Use()
	rtr.Use(middleware.SetCommonHeaders())
	// Configure Gin to use Zap
	rtr.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/api/v1/healthcheck"},
	}))
	rtr.Use(gin.Recovery())
	return rtr
}

// Start starts the HTTP server in a goroutine
func (a *APIServer) Start() error {
	// Start API
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal("Server listen and serve error", zap.Error(err))
			os.Exit(1)
		}
	}()
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
//	@version		0.5
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
	defer func() { _ = logger.Sync() }()

	// Loading APIServer config
	cfg, err := config.LoadAPIServerConfig()
	if err != nil {
		logger.Fatal("Error loading APIServer config", zap.Error(err))
	}
	logger.Info("Configuration loaded successfully")

	// Initialize database connection
	// TODO: Define a better context
	dbClient, err := dbclient.NewDBClient(cfg.DBURL, logger, context.TODO())
	if err != nil {
		logger.Fatal("Could not establish database connection", zap.Error(err))
	}
	defer dbClient.Close()

	// Initializing gRPC AgentClient
	agentClient, err := clients.NewGRPCAgentClient(cfg.AgentURL, logger)
	if err != nil {
		logger.Fatal("Failed to create gRPC client", zap.Error(err))
	}

	// Initializing repositories
	accountRepo := repositories.NewAccountRepository(dbClient)
	clusterRepo := repositories.NewClusterRepository(dbClient)
	instanceRepo := repositories.NewInstanceRepository(dbClient)
	expenseRepo := repositories.NewExpenseRepository(dbClient)
	eventRepo := repositories.NewEventRepository(dbClient)
	actionRepo := repositories.NewActionRepository(dbClient)

	// Initializing services
	accountService := services.NewAccountService(accountRepo)
	clusterServiceOpts := services.ClusterServiceOptions{
		AgentRequestTimeout: cfg.AgentRequestTimeout,
	}
	clusterService := services.NewClusterService(clusterRepo, agentClient, clusterServiceOpts)
	instanceService := services.NewInstanceService(instanceRepo)
	expenseService := services.NewExpenseService(expenseRepo)
	eventService := services.NewEventService(eventRepo)
	actionService := services.NewActionService(actionRepo)
	overviewService := services.NewOverviewService(clusterRepo, instanceRepo, accountRepo)

	// Initializing handlers
	handlers := APIHandlers{
		AccountHandler:     handlers.NewAccountHandler(accountService, logger),
		ClusterHandler:     handlers.NewClusterHandler(clusterService, logger),
		InstanceHandler:    handlers.NewInstanceHandler(instanceService, logger),
		ExpenseHandler:     handlers.NewExpenseHandler(expenseService, logger),
		EventHandler:       handlers.NewEventHandler(eventService, logger),
		ActionHandler:      handlers.NewActionHandler(actionService, logger),
		OverviewHandler:    handlers.NewOverviewHandler(overviewService, logger),
		HealthCheckHandler: handlers.NewHealthCheckHandler(dbClient, logger),
	}

	// Setup router
	engine := setupGin(logger)
	Setup(engine, handlers)

	// parsing cfg.DBURL for removing user:password when logging
	dburl, _ := url.Parse(cfg.DBURL)
	dbHost := dburl.Hostname() + ":" + dburl.Port()

	logger.Info("ClusterIQ API server started",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("listenURL", cfg.ListenURL),
		zap.String("dbURL", dbHost),
		zap.String("agentURL", cfg.AgentURL),
	)

	// Initializing APIServer instance
	api := NewAPIServer(cfg.ListenURL, logger, engine)
	if err := api.Run(); err != nil {
		logger.Fatal("Error running API server", zap.Error(err))
	}
}
