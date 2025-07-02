package router

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	ClusterHandler *handlers.ClusterHandler
}

func Setup(engine *gin.Engine, deps Dependencies) {
	baseGroup := engine.Group("/api/v1")
	{
		setupClustersRoutes(baseGroup, deps.ClusterHandler)
		//setupHealthcheckRoutes(baseGroup) // TODO: add dependenctis for DB
		//setupScheduledActionsRoutes(baseGroup)
		//setupExpensesRoutes(baseGroup)
		//setupInstancesRoutes(baseGroup)
		//setupAccountsRoutes(baseGroup)
		//setupEventsRoutes(baseGroup)
		//setupOverviewRoutes(baseGroup)
		//setupInventoryRoutes(baseGroup)
	}
}

func setupClustersRoutes(baseGroup *gin.RouterGroup, handler *handlers.ClusterHandler) {
	clustersGroup := baseGroup.Group("/clusters")
	clustersGroup.GET("", handler.ListClusters)
	//clustersGroup.GET("/:cluster_id", r.api.HandlerGetClustersByID)
	//clustersGroup.GET("/:cluster_id/instances", r.api.HandlerGetInstancesOnCluster)
	//clustersGroup.GET("/:cluster_id/tags", r.api.HandlerGetClusterTags)
	//clustersGroup.GET("/:cluster_id/events", r.api.HandlerGetClusterEvents)
	//clustersGroup.POST("", r.api.HandlerPostCluster, r.api.HandlerRefreshInventory)
	//clustersGroup.POST("/:cluster_id/power_on", r.api.HandlerPowerOnCluster)
	//clustersGroup.POST("/:cluster_id/power_off", r.api.HandlerPowerOffCluster)
	//clustersGroup.DELETE("/:cluster_id", r.api.HandlerDeleteCluster)
	//clustersGroup.PATCH("/:cluster_id", r.api.HandlerPatchCluster)
}
