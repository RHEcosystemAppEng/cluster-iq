package main

import (
	"fmt"
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ==================== Inventory  Handlers ====================

// HandlerRefreshInventory handles the request for refreshing the entire
// inventory just after a full scan. This method is used for recalculating some
// values and mark the missing clusters as "terminated"
//
//	@Summary		Refresh data on inventory
//	@Description	Recalculating some values and mark the missing clusters as "terminated"
//	@Tags			Inventory
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		500	{object}	nil
//	@Router			/inventory/refresh [post]
func (a APIServer) HandlerRefreshInventory(c *gin.Context) {
	if err := a.sql.RefreshInventory(); err != nil {
		a.logger.Error("Can't refresh inventory data on DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}
	// This function doesn't return any 200OK code for preventing duplicated responses
}

// HandlerGetSystemEvents handles the request for obtain the list of system events
//
//	@Summary		Obtain system events
//	@Description	Returns a list of events
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	EventsListResponse
//	@Failure		500	{object}	nil
//	@Router			/events [get]
func (a APIServer) HandlerGetSystemEvents(c *gin.Context) {
	a.logger.Debug("Retrieving system-wide events")

	dbEvents, err := a.sql.GetSystemEvents()
	if err != nil {
		a.logger.Error("Failed to retrieve system-wide events", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse("failed to retrieve system-wide events"))
		return
	}

	appEvents := events.ToSystemAuditEvents(dbEvents)
	c.PureJSON(http.StatusOK, NewSystemEventsListResponse(appEvents))
}

// HandlerGetClusterEvents handles the request for obtain the list of events of a Cluster
//
//	@Summary		Obtain cluster events
//	@Description	Returns a list of events belonging to a cluster given by ID
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	EventsListResponse
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id}/events [get]
func (a APIServer) HandlerGetClusterEvents(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Retrieving cluster events", zap.String("cluster_id", clusterID))

	dbEvents, err := a.sql.GetClusterEvents(clusterID)
	if err != nil {
		a.logger.Error("Failed to retrieve cluster events",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		c.PureJSON(http.StatusInternalServerError,
			NewGenericErrorResponse("failed to retrieve cluster events"))
		return
	}
	appEvents := events.ToAuditEvents(dbEvents)
	c.PureJSON(http.StatusOK, NewEventsListResponse(appEvents))
}

// HandlerGetInventoryOverview handles the request to obtain an overview of the inventory
//
//	@Summary		Obtain an inventory overview
//	@Description	Returns an overview of the inventory
//	@Tags			Overview
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	models.OverviewSummary
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/overview	[get]
func (a APIServer) HandlerGetInventoryOverview(c *gin.Context) {
	a.logger.Debug("Retrieving overview data")

	overview, err := a.getInventoryOverview()
	if err != nil {
		a.logger.Error("Error when generating Inventory Overview data", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse("failed to retrieve inventory overview"))
		return
	}
	c.PureJSON(http.StatusOK, overview)
}

// getInventoryOverview retrieves all components of the inventory overview.
func (a APIServer) getInventoryOverview() (models.OverviewSummary, error) {
	var overview models.OverviewSummary

	// Get clusters summary
	clusters, err := a.sql.GetClustersOverview()
	if err != nil {
		return models.OverviewSummary{}, fmt.Errorf("failed to get clusters overview: %w", err)
	}
	overview.Clusters = clusters

	// Get instances summary
	instances, err := a.sql.GetInstancesOverview()
	if err != nil {
		return models.OverviewSummary{}, fmt.Errorf("failed to get instances overview: %w", err)
	}
	overview.Instances = instances

	// Get providers summary
	providers, err := a.sql.GetProvidersOverview()
	if err != nil {
		return models.OverviewSummary{}, fmt.Errorf("failed to get providers overview: %w", err)
	}
	overview.Providers = providers

	// Get scanner last scan timestamp
	scannerLastScan, err := a.sql.GetScannerLastScanTimestamp()
	if err != nil {
		return models.OverviewSummary{}, fmt.Errorf("failed to get scanner last scan timestamp: %w", err)
	}
	overview.Scanner = models.Scanner{LastScanTimestamp: scannerLastScan}

	return overview, nil
}
