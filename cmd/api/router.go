package main

import (
	"github.com/RHEcosystemAppEng/cluster-iq/cmd/api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Router struct {
	engine *gin.Engine
	api    *APIServer
}

func NewRouter(api *APIServer) *Router {
	return &Router{
		engine: api.router,
		api:    api,
	}
}

func (r *Router) SetupRoutes() {
	// API Endpoints
	r.setupSwagger()
	baseGroup := r.engine.Group("/api/v1")
	r.setupHealthcheckRoutes(baseGroup)
	r.setupExpensesGroupRoutes(baseGroup)
	r.setupInstancesRoutes(baseGroup)
	r.setupClustersRoutes(baseGroup)
	r.setupAccountsRoutes(baseGroup)
	r.setupSwaggerRoutes(baseGroup)
}

func (r *Router) setupHealthcheckRoutes(baseGroup *gin.RouterGroup) {
	healthcheckGroup := baseGroup.Group("/healthcheck")
	healthcheckGroup.GET("", r.api.HandlerHealthCheck)
}

func (r *Router) setupExpensesGroupRoutes(baseGroup *gin.RouterGroup) {
	expensesGroup := baseGroup.Group("/expenses")
	expensesGroup.GET("", r.api.HandlerGetExpenses)
	expensesGroup.GET("/:instance_id", r.api.HandlerGetExpensesByInstance)
	expensesGroup.POST("", r.api.HandlerPostExpense)
}

func (r *Router) setupInstancesRoutes(baseGroup *gin.RouterGroup) {
	instancesGroup := baseGroup.Group("/instances")
	instancesGroup.Use(r.api.HandlerRefreshInventory)
	instancesGroup.GET("", r.api.HandlerGetInstances)
	instancesGroup.GET("/expense_update", r.api.HandlerGetInstancesForBillingUpdate)
	instancesGroup.GET("/:instance_id", r.api.HandlerGetInstanceByID)
	instancesGroup.POST("", r.api.HandlerPostInstance)
	instancesGroup.DELETE("/:instance_id", r.api.HandlerDeleteInstance)
	instancesGroup.PATCH("/:instance_id", r.api.HandlerPatchInstance)
}

func (r *Router) setupClustersRoutes(baseGroup *gin.RouterGroup) {
	clustersGroup := baseGroup.Group("/clusters")
	clustersGroup.Use(r.api.HandlerRefreshInventory)
	clustersGroup.GET("", r.api.HandlerGetClusters)
	clustersGroup.GET("/:cluster_id", r.api.HandlerGetClustersByID)
	clustersGroup.GET("/:cluster_id/instances", r.api.HandlerGetInstancesOnCluster)
	clustersGroup.GET("/:cluster_id/tags", r.api.HandlerGetClusterTags)
	clustersGroup.POST("", r.api.HandlerPostCluster)
	clustersGroup.POST("/:cluster_id/power_on", r.api.HandlerPowerOnCluster)
	clustersGroup.POST("/:cluster_id/power_off", r.api.HandlerPowerOffCluster)
	clustersGroup.DELETE("/:cluster_id", r.api.HandlerDeleteCluster)
	clustersGroup.PATCH("/:cluster_id", r.api.HandlerPatchCluster)
}

func (r *Router) setupAccountsRoutes(baseGroup *gin.RouterGroup) {
	accountsGroup := baseGroup.Group("/accounts")
	accountsGroup.GET("", r.api.HandlerGetAccounts)
	accountsGroup.GET("/:account_name", r.api.HandlerGetAccountsByName)
	accountsGroup.GET("/:account_name/clusters", r.api.HandlerGetClustersOnAccount)
	accountsGroup.POST("", r.api.HandlerPostAccount)
	accountsGroup.DELETE("/:account_name", r.api.HandlerDeleteAccount)
	accountsGroup.PATCH("/:account_name", r.api.HandlerPatchAccount)
}

func (r *Router) setupSwaggerRoutes(baseGroup *gin.RouterGroup) {
	baseGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *Router) setupSwagger() {
	docs.SwaggerInfo.Title = "Cluster IP API doc"
	docs.SwaggerInfo.Description = "This the API of the ClusterIQ project"
	docs.SwaggerInfo.Version = "0.3"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
}
