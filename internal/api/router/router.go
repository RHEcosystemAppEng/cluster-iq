package router

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// Dependencies holds all the handlers that the router needs.
type Dependencies struct {
	AccountHandler     *handlers.AccountHandler
	ClusterHandler     *handlers.ClusterHandler
	InstanceHandler    *handlers.InstanceHandler
	ExpenseHandler     *handlers.ExpenseHandler
	EventHandler       *handlers.EventHandler
	ActionHandler      *handlers.ActionHandler
	HealthCheckHandler *handlers.HealthCheckHandler
}

// Setup configures the API routes.
func Setup(engine *gin.Engine, deps Dependencies) {
	baseGroup := engine.Group("/api/v1")
	{
		setupHealthCheckRoutes(baseGroup, deps.HealthCheckHandler)
		setupAccountRoutes(baseGroup, deps.AccountHandler)
		setupClusterRoutes(baseGroup, deps.ClusterHandler, deps.InstanceHandler, deps.EventHandler)
		setupInstanceRoutes(baseGroup, deps.InstanceHandler)
		setupExpenseRoutes(baseGroup, deps.ExpenseHandler)
		setupEventRoutes(baseGroup, deps.EventHandler)
		setupActionRoutes(baseGroup, deps.ActionHandler)
	}
}

func setupHealthCheckRoutes(group *gin.RouterGroup, handler *handlers.HealthCheckHandler) {
	group.GET("/health", handler.Check)
}

func setupAccountRoutes(group *gin.RouterGroup, handler *handlers.AccountHandler) {
	accounts := group.Group("/accounts")
	{
		accounts.GET("", handler.List)
		accounts.POST("", handler.Create)
		accounts.GET("/:name", handler.GetByName)
		accounts.DELETE("/:name", handler.Delete)
	}
}

func setupClusterRoutes(group *gin.RouterGroup, clusterHandler *handlers.ClusterHandler, instanceHandler *handlers.InstanceHandler, eventHandler *handlers.EventHandler) {
	clusters := group.Group("/clusters")
	{
		clusters.GET("", clusterHandler.List)
		// clusters.GET("/summary", clusterHandler.GetSummary)
		clusters.POST("", clusterHandler.Create)
		clusters.GET("/:id", clusterHandler.Get)
		clusters.DELETE("/:id", clusterHandler.Delete)
		clusters.PATCH("/:id", clusterHandler.Update)
		clusters.POST("/:id/power_on", clusterHandler.PowerOn)
		clusters.POST("/:id/power_off", clusterHandler.PowerOff)
		clusters.GET("/:id/tags", clusterHandler.GetTags)
		clusters.GET("/:id/instances", instanceHandler.ListByCluster)
		clusters.GET("/:id/events", eventHandler.ListByCluster)
	}
}

func setupInstanceRoutes(group *gin.RouterGroup, handler *handlers.InstanceHandler) {
	instances := group.Group("/instances")
	{
		instances.GET("", handler.List)
		// instances.GET("/summary", handler.GetSummary)
		instances.GET("/:id", handler.Get)
	}
}

func setupExpenseRoutes(group *gin.RouterGroup, handler *handlers.ExpenseHandler) {
	expenses := group.Group("/expenses")
	{
		expenses.GET("", handler.List)
	}
}

func setupEventRoutes(group *gin.RouterGroup, handler *handlers.EventHandler) {
	events := group.Group("/events")
	{
		events.GET("/system", handler.ListSystem)
	}
}

func setupActionRoutes(group *gin.RouterGroup, handler *handlers.ActionHandler) {
	actions := group.Group("/actions/scheduled")
	{
		// actions.GET("", handler.ListScheduled)
		// actions.GET("/:id", handler.GetScheduled)
		actions.PATCH("/:id/enable", handler.EnableScheduled)
		actions.PATCH("/:id/disable", handler.DisableScheduled)
	}
}
