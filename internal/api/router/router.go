package router

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// APIHandlers holds all the handlers that the router needs.
type APIHandlers struct {
	AccountHandler     *handlers.AccountHandler
	ClusterHandler     *handlers.ClusterHandler
	InstanceHandler    *handlers.InstanceHandler
	ExpenseHandler     *handlers.ExpenseHandler
	EventHandler       *handlers.EventHandler
	ActionHandler      *handlers.ActionHandler
	OverviewHandler    *handlers.OverviewHandler
	HealthCheckHandler *handlers.HealthCheckHandler
}

// Setup configures the API routes.
func Setup(engine *gin.Engine, handlers APIHandlers) {
	baseGroup := engine.Group("/api/v1")
	{
		setupHealthCheckRoutes(baseGroup, handlers.HealthCheckHandler)
		setupAccountRoutes(baseGroup, handlers.AccountHandler)
		setupClusterRoutes(baseGroup, handlers.ClusterHandler, handlers.InstanceHandler, handlers.EventHandler)
		setupInstanceRoutes(baseGroup, handlers.InstanceHandler)
		setupExpenseRoutes(baseGroup, handlers.ExpenseHandler)
		setupEventRoutes(baseGroup, handlers.EventHandler)
		setupActionRoutes(baseGroup, handlers.ActionHandler)
		setupOverviewRoutes(baseGroup, handlers.OverviewHandler)
	}
}

func setupHealthCheckRoutes(group *gin.RouterGroup, handler *handlers.HealthCheckHandler) {
	group.GET("/healthcheck", handler.Check)
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
		clusters.POST("", clusterHandler.Create)
		clusters.GET("/:id", clusterHandler.Get)
		clusters.DELETE("/:id", clusterHandler.Delete)
		clusters.PATCH("/:id", clusterHandler.Update)
		clusters.POST("/:id/power_on", clusterHandler.PowerOn)
		clusters.POST("/:id/power_off", clusterHandler.PowerOff)
		clusters.GET("/:id/tags", clusterHandler.GetTags)
		clusters.GET("/:id/events", eventHandler.ListByCluster)
	}
}

func setupInstanceRoutes(group *gin.RouterGroup, handler *handlers.InstanceHandler) {
	instances := group.Group("/instances")
	{
		instances.GET("", handler.List)
		instances.POST("", handler.Create)
		instances.GET("/:id", handler.Get)
	}
}

func setupExpenseRoutes(group *gin.RouterGroup, handler *handlers.ExpenseHandler) {
	expenses := group.Group("/expenses")
	{
		expenses.GET("", handler.List)
		expenses.POST("", handler.Create)
	}
}

func setupEventRoutes(group *gin.RouterGroup, handler *handlers.EventHandler) {
	events := group.Group("/events")
	{
		events.GET("/system", handler.ListSystem)
	}
}

func setupActionRoutes(group *gin.RouterGroup, handler *handlers.ActionHandler) {
	schedule := group.Group("/schedule")
	{
		schedule.GET("", handler.ListScheduled)
		schedule.GET("/:id", handler.GetScheduled)
		schedule.PATCH("/:id/enable", handler.EnableScheduled)
		schedule.PATCH("/:id/disable", handler.DisableScheduled)
	}
}

func setupOverviewRoutes(group *gin.RouterGroup, handler *handlers.OverviewHandler) {
	overview := group.Group("/overview")
	{
		overview.GET("", handler.Get)
	}
}
