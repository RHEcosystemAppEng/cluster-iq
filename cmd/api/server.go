// Package main is the entry point for the ClusterIQ API server.
// It initializes the API server, sets up routes, and handles server lifecycle events.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/cmd/api/docs"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/config"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
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
func NewAPIServer(cfg *config.APIServerConfig, logger *zap.Logger) *APIServer {
	// Configuring GIN router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Configure GIN to use ZAP
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Configure HTTP server
	server := &http.Server{
		Addr:    cfg.ListenURL,
		Handler: router.Handler(),
	}

	// Creating gRPC client
	gRPCClient, err := NewAPIGRPCClient(cfg.AgentURL, logger)
	if err != nil {
		logger.Error("Cannot create gRPC Client", zap.String("agent_url", cfg.AgentURL), zap.Error(err))
		return nil
	}

	// Creating DB client
	sqlClient, err := NewAPISQLClient(cfg.DBURL, logger)
	if err != nil {
		logger.Error("Cannot create SQL Client", zap.String("db_url", cfg.DBURL), zap.Error(err))
		return nil
	}

	return &APIServer{
		cfg:    cfg,
		logger: logger,
		server: server,
		router: router,
		grpc:   gRPCClient,
		sql:    sqlClient,
	}
}

func init() {
}

// signalHandler handles OS signals for graceful server shutdown.  It shuts
// down the server when a SIGTERM signal is received. This function was included
// for better integration on K8s/OCP
//
// Parameters:
// - signal: The OS signal to handle.
func (a APIServer) signalHandler(signal os.Signal) {
	if signal == syscall.SIGTERM {
		ctx, cancel := context.WithTimeout(context.Background(), APITimeoutSeconds*time.Second)
		defer cancel()
		a.logger.Warn("SIGTERM signal received. Stopping ClusterIQ API server")

		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Fatal("API Shutdown error", zap.Error(err))
			os.Exit(-1)
		}
	} else {
		a.logger.Warn("Ignoring signal: ", zap.String("signal_id", signal.String()))
	}
}

// addHeaders adds the required HTTP headers for API working
func addHeaders(c *gin.Context) {
	// To deal with CORS
	c.Header("Access-Control-Allow-Origin", "*")
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
	api := NewAPIServer(cfg, logger)

	api.logger.Info("==================== Starting ClusterIQ API ====================",
		zap.String("version", version),
		zap.String("commit", commit),
	)
	api.logger.Info("Connection properties", zap.String("api_url", api.cfg.ListenURL), zap.String("db_url", api.cfg.DBURL), zap.String("agent_url", api.cfg.AgentURL))
	api.logger.Debug("Debug Mode active!")

	// Preparing API Endpoints
	baseGroup := api.router.Group("/api/v1")
	{
		healthcheckGroup := baseGroup.Group("/healthcheck")
		{
			healthcheckGroup.GET("", api.HandlerHealthCheck)
		}
		expensesGroup := baseGroup.Group("/expenses")
		{
			expensesGroup.GET("", api.HandlerGetExpenses)
			expensesGroup.GET("/:instance_id", api.HandlerGetExpensesByInstance)
			expensesGroup.POST("", api.HandlerPostExpense)
		}
		instancesGroup := baseGroup.Group("/instances")
		instancesGroup.Use(api.HandlerRefreshInventory)
		{
			instancesGroup.GET("", api.HandlerGetInstances)
			instancesGroup.GET("/expense_update", api.HandlerGetInstancesForBillingUpdate)
			instancesGroup.GET("/:instance_id", api.HandlerGetInstanceByID)
			instancesGroup.POST("", api.HandlerPostInstance)
			instancesGroup.DELETE("/:instance_id", api.HandlerDeleteInstance)
			instancesGroup.PATCH("/:instance_id", api.HandlerPatchInstance)
		}

		clustersGroup := baseGroup.Group("/clusters")
		clustersGroup.Use(api.HandlerRefreshInventory)
		{
			clustersGroup.GET("", api.HandlerGetClusters)
			clustersGroup.GET("/:cluster_id", api.HandlerGetClustersByID)
			clustersGroup.GET("/:cluster_id/instances", api.HandlerGetInstancesOnCluster)
			clustersGroup.GET("/:cluster_id/tags", api.HandlerGetClusterTags)
			clustersGroup.POST("", api.HandlerPostCluster)
			clustersGroup.POST("/:cluster_id/power_on", api.HandlerPowerOnCluster)
			clustersGroup.POST("/:cluster_id/power_off", api.HandlerPowerOffCluster)
			clustersGroup.DELETE("/:cluster_id", api.HandlerDeleteCluster)
			clustersGroup.PATCH("/:cluster_id", api.HandlerPatchCluster)
		}

		accountsGroup := baseGroup.Group("/accounts")
		{
			accountsGroup.GET("", api.HandlerGetAccounts)
			accountsGroup.GET("/:account_name", api.HandlerGetAccountsByName)
			accountsGroup.GET("/:account_name/clusters", api.HandlerGetClustersOnAccount)
			accountsGroup.POST("", api.HandlerPostAccount)
			accountsGroup.DELETE("/:account_name", api.HandlerDeleteAccount)
			accountsGroup.PATCH("/:account_name", api.HandlerPatchAccount)
		}
		baseGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Swagger endpoint
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Cluster IP API doc"
	docs.SwaggerInfo.Description = "This the API of the ClusterIQ project"
	docs.SwaggerInfo.Version = "0.3"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// Start API
	go func() {
		if err := api.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Server listen and serve error", zap.Error(err))
			os.Exit(-1)
		}
	}()

	// Log API is running
	api.logger.Info("API Ready to serve")

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGTERM)
	s := <-quitChan
	api.signalHandler(s)
	api.logger.Info("API server stopped")
}
