// TODO: Read this article https://ferencfbin.medium.com/golang-own-structscan-method-for-sql-rows-978c5c80f9b5
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/RHEcosystemAppEng/cluster-iq/cmd/api/docs"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	ciqLogger "github.com/RHEcosystemAppEng/cluster-iq/internal/logger"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	// swagger embed files
	// swagger embed files
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var (
	// version reflects the current version of the API
	version string
	// commit reflects the git short-hash of the compiled version
	commit string
	// TODO: comment rest of global vars
	inven  *inventory.Inventory
	router *gin.Engine
	apiURL string
	dbURL  string
	dbPass string
	logger *zap.Logger
	db     *sqlx.DB
)

const connStr = "postgresql://user:password@localhost:5432/clusteriq?sslmode=disable"

func init() {
	// Logging config
	logger = ciqLogger.NewLogger()

	// Getting config
	apiHost := os.Getenv("CIQ_API_HOST")
	apiPort := os.Getenv("CIQ_API_PORT")
	dbHost := os.Getenv("CIQ_DB_HOST")
	dbPort := os.Getenv("CIQ_DB_PORT")
	dbPass = os.Getenv("CIQ_DB_PASS")
	apiURL = fmt.Sprintf("%s:%s", apiHost, apiPort)
	dbURL = fmt.Sprintf("%s:%s", dbHost, dbPort)

	// Initializaion global vars
	inven = inventory.NewInventory()
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()

	// Configure GIN to use ZAP
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
}

func addHeaders(c *gin.Context) {
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
		instancesGroup := router.Group("/instances")
		{
			instancesGroup.GET("/", HandlerGetInstances)
			instancesGroup.GET("/:instance_id", HandlerGetInstancesByID)
			instancesGroup.POST("/", HandlerPostInstance)
			instancesGroup.DELETE("/", HandlerDeleteInstance)
			instancesGroup.PATCH("/", HandlerPatchInstance)
		}

		clustersGroup := router.Group("/clusters")
		{
			clustersGroup.GET("/", HandlerGetClusters)
			clustersGroup.GET("/:cluster_name", HandlerGetClustersByName)
			clustersGroup.GET("/:cluster_name/instances", HandlerGetInstancesOnCluster)
			clustersGroup.POST("/", HandlerPostCluster)
			clustersGroup.DELETE("/", HandlerDeleteCluster)
			clustersGroup.PATCH("/", HandlerPatchCluster)
		}

		accountsGroup := router.Group("/accounts")
		{
			accountsGroup.GET("/", HandlerGetAccounts)
			accountsGroup.GET("/:account_name", HandlerGetAccountsByName)
			accountsGroup.GET("/:account_name/clusters", HandlerGetClustersOnAccount)
			accountsGroup.POST("/", HandlerPostAccount)
			accountsGroup.DELETE("/", HandlerDeleteAccount)
			accountsGroup.PATCH("/", HandlerPatchAccount)
		}
		baseGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Swagger endpoint
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Cluster IP API doc"
	docs.SwaggerInfo.Description = "This the API of the ClusterIQ project"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// PGSQL connection
	var err error
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		logger.Error("Can't connect to PSQL DB", zap.Error(err))
	}

	// Start API
	logger.Info("API Ready to serve")
	router.Run(apiURL)
}
