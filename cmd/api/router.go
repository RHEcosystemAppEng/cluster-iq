package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func (a *APIServer) setupRouter() {
	// API Endpoints
	baseGroup := a.router.Group("/api/v1")
	a.setupHealthcheckRoutes(baseGroup)
	a.setupExpensesGroupRoutes(baseGroup)
	a.setupInstancesRoutes(baseGroup)
	a.setupClustersRoutes(baseGroup)
	a.setupAccountsRoutes(baseGroup)
	a.setupSwaggerRoutes(baseGroup)
}

func (a *APIServer) setupHealthcheckRoutes(baseGroup *gin.RouterGroup) {
	healthcheckGroup := baseGroup.Group("/healthcheck")
	healthcheckGroup.GET("", a.HandlerHealthCheck)
}

func (a *APIServer) setupExpensesGroupRoutes(baseGroup *gin.RouterGroup) {
	expensesGroup := baseGroup.Group("/expenses")
	expensesGroup.GET("", a.HandlerGetExpenses)
	expensesGroup.GET("/:instance_id", a.HandlerGetExpensesByInstance)
	expensesGroup.POST("", a.HandlerPostExpense)
}

func (a *APIServer) setupInstancesRoutes(baseGroup *gin.RouterGroup) {
	instancesGroup := baseGroup.Group("/instances")
	instancesGroup.Use(a.HandlerRefreshInventory)
	instancesGroup.GET("", a.HandlerGetInstances)
	instancesGroup.GET("/expense_update", a.HandlerGetInstancesForBillingUpdate)
	instancesGroup.GET("/:instance_id", a.HandlerGetInstanceByID)
	instancesGroup.POST("", a.HandlerPostInstance)
	instancesGroup.DELETE("/:instance_id", a.HandlerDeleteInstance)
	instancesGroup.PATCH("/:instance_id", a.HandlerPatchInstance)
}

func (a *APIServer) setupClustersRoutes(baseGroup *gin.RouterGroup) {
	clustersGroup := baseGroup.Group("/clusters")
	clustersGroup.Use(a.HandlerRefreshInventory)
	clustersGroup.GET("", a.HandlerGetClusters)
	clustersGroup.GET("/:cluster_id", a.HandlerGetClustersByID)
	clustersGroup.GET("/:cluster_id/instances", a.HandlerGetInstancesOnCluster)
	clustersGroup.GET("/:cluster_id/tags", a.HandlerGetClusterTags)
	clustersGroup.POST("", a.HandlerPostCluster)
	clustersGroup.POST("/:cluster_id/power_on", a.HandlerPowerOnCluster)
	clustersGroup.POST("/:cluster_id/power_off", a.HandlerPowerOffCluster)
	clustersGroup.DELETE("/:cluster_id", a.HandlerDeleteCluster)
	clustersGroup.PATCH("/:cluster_id", a.HandlerPatchCluster)
}

func (a *APIServer) setupAccountsRoutes(baseGroup *gin.RouterGroup) {
	accountsGroup := baseGroup.Group("/accounts")
	accountsGroup.GET("", a.HandlerGetAccounts)
	accountsGroup.GET("/:account_name", a.HandlerGetAccountsByName)
	accountsGroup.GET("/:account_name/clusters", a.HandlerGetClustersOnAccount)
	accountsGroup.POST("", a.HandlerPostAccount)
	accountsGroup.DELETE("/:account_name", a.HandlerDeleteAccount)
	accountsGroup.PATCH("/:account_name", a.HandlerPatchAccount)
}

func (a *APIServer) setupSwaggerRoutes(baseGroup *gin.RouterGroup) {
	baseGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
