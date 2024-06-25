package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/cmd/api/docs"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

const (
	APITimeoutSeconds = 10
)

var (
	// version reflects the current version of the API
	version string
	// commit reflects the git short-hash of the compiled version
	commit string
	// TODO: comment rest of global vars
	inven *inventory.Inventory
	// Gin Router for serving API
	router *gin.Engine
	// Gin Server
	server *http.Server
	// API serving URL
	apiURL string
	// DB URL (example: "postgresql://user:password@pgsql:5432/clusteriq?sslmode=disable")
	dbURL string
	// Loggin verbosity level. For more info, check zap log library
	logLevel string
	// Logger object
	logger *zap.Logger
	// DB object to manage DB connections
	db *sqlx.DB
)

func init() {
	// Logger creation
	logger = ciqLogger.NewLogger()

	// Getting Env Vars for config
	apiURL = os.Getenv("CIQ_API_LISTEN_URL")
	dbURL = os.Getenv("CIQ_DB_URL")

	// Initializaion global vars
	inven = inventory.NewInventory()
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// Configure GIN to use ZAP
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
}

// signalHandler for managing incoming OS signals
func signalHandler(signal os.Signal) {
	if signal == syscall.SIGTERM {
		ctx, cancel := context.WithTimeout(context.Background(), APITimeoutSeconds*time.Second)
		defer cancel()
		logger.Warn("SIGTERM signal received. Stopping ClusterIQ API server")
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("API Shutdown error", zap.Error(err))
			os.Exit(-1)
		}
	} else {
		logger.Warn("Ignoring signal: ", zap.String("signal_id", signal.String()))
	}
}

// addHeaders adds the requerired HTTP headers for API working
func addHeaders(c *gin.Context) {
	// To deal with CORS
	c.Header("Access-Control-Allow-Origin", "*")
}

//	@title			ClusterIQ API
//	@version		1.0
//	@description	This is the API of the ClusterIQ cloud inventory software
//	@tersOfService	http://swagger.io/ters/

//	@contact.name	ClusterIQ Team
//	@contact.email	vbelouso@redhat.com nnaamneh@redhat.com avillega@redhat.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// Ignore Logger sync error
	defer func() { _ = logger.Sync() }()

	logger.Info("Starting ClusterIQ API", zap.String("version", version), zap.String("commit", commit))
	logger.Info("Connection properties", zap.String("api_url", apiURL), zap.String("db_url", dbURL))
	logger.Debug("Debug Mode active!")

	// Preparing API Endpoints
	baseGroup := router.Group("/api/v1")
	{
		healthcheckGroup := baseGroup.Group("/healthcheck")
		{
			healthcheckGroup.GET("", HandlerHealthCheck)
		}

		instancesGroup := baseGroup.Group("/instances")
		{
			instancesGroup.GET("", HandlerGetInstances)
			instancesGroup.GET("/:instance_id", HandlerGetInstanceByID)
			instancesGroup.POST("", HandlerPostInstance)
			instancesGroup.DELETE("/:instance_id", HandlerDeleteInstance)
			instancesGroup.PATCH("/:instance_id", HandlerPatchInstance)
		}

		clustersGroup := baseGroup.Group("/clusters")
		{
			clustersGroup.GET("", HandlerGetClusters)
			clustersGroup.GET("/:cluster_id", HandlerGetClustersByID)
			clustersGroup.GET("/:cluster_id/instances", HandlerGetInstancesOnCluster)
			clustersGroup.GET("/:cluster_id/tags", HandlerGetClusterTags)
			clustersGroup.POST("", HandlerPostCluster)
			clustersGroup.DELETE("/:cluster_id", HandlerDeleteCluster)
			clustersGroup.PATCH("/:cluster_id", HandlerPatchCluster)
		}

		accountsGroup := baseGroup.Group("/accounts")
		{
			accountsGroup.GET("", HandlerGetAccounts)
			accountsGroup.GET("/:account_name", HandlerGetAccountsByName)
			accountsGroup.GET("/:account_name/clusters", HandlerGetClustersOnAccount)
			accountsGroup.POST("", HandlerPostAccount)
			accountsGroup.DELETE("/:account_name", HandlerDeleteAccount)
			accountsGroup.PATCH("/:account_name", HandlerPatchAccount)
		}
		baseGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Swagger endpoint
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Cluster IP API doc"
	docs.SwaggerInfo.Description = "This the API of the ClusterIQ project"
	docs.SwaggerInfo.Version = "0.2"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// PGSQL connection
	var err error
	db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		logger.Error("Can't connect to PSQL DB", zap.Error(err))
		return
	}

	server = &http.Server{
		Addr:    apiURL,
		Handler: router.Handler(),
	}

	// Start API
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server listen and serve error", zap.Error(err))
			os.Exit(-1)
		}
	}()

	logger.Info("API Ready to serve")

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGTERM)
	s := <-quitChan
	signalHandler(s)
	logger.Info("API server stopped")
}
