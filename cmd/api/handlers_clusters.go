package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ==================== Clusters      Handlers ====================

// HandlerGetClusters handles the request for obtaining the entire Cluster list
//
//	@Summary		Obtain every Cluster
//	@Description	Returns a list of Clusters with a single instance filtered by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ClusterListResponse
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/clusters [get]
func (a APIServer) HandlerGetClusters(c *gin.Context) {
	a.logger.Debug("Retrieving complete Clusters inventory")

	clusters, err := a.sql.GetClusters()
	if err != nil {
		a.logger.Error("Can't retrieve Clusters list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	var response []dto.ClusterDTOResponse
	for _, cluster := range clusters {
		response = append(response, *cluster.ToClusterDTOResponse())
	}

	c.PureJSON(http.StatusOK, dto.NewClusterDTOResponseList(response))
}

// TODO: After reviewing the gRPC structure, check if this method still being useful
// HandlerGetClustersByID handles the request for obtain a Cluster by its Name
//
//	@Summary		Obtain a single Cluster by its Name
//	@Description	Returns a list of Clusters with a single Cluster filtered by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	ClusterListResponse
//	@Failure		404			{object}	nil
//	@Router			/clusters/{cluster_id} [get]
func (a APIServer) HandlerGetClustersByID(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Retrieving Cluster Tags by ID", zap.String("cluster_id", clusterID))

	clusters, err := a.sql.GetClusterByID(clusterID)
	if err != nil {
		a.logger.Error("Cluster not found", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	c.PureJSON(http.StatusOK, dto.NewClusterDTOResponseList([]dto.ClusterDTOResponse{*clusters.ToClusterDTOResponse()}))
}

// HandlerGetInstancesOnCluster handles the request for obtain the list of Instances belonging to a specific Cluster
//
//	@Summary		Obtain Instances list belonging to a Cluster
//	@Description	Returns a list of Instances belonging to a Cluster given by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	InstanceListResponse
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id}/instances [get]
func (a APIServer) HandlerGetInstancesOnCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Retrieving Cluster's Instances", zap.String("cluster_id", clusterID))

	instances, err := a.sql.GetInstancesOnCluster(clusterID)
	if err != nil {
		a.logger.Error("Can't retrieve instances on cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Transforming into DTO type
	var response []dto.InstanceDTOResponse
	for _, cluster := range instances {
		response = append(response, *cluster.ToInstanceDTOResponse())
	}

	c.PureJSON(http.StatusOK, dto.NewInstanceDTOResponseList(response))
}

// HandlerGetClusterTags handles the request for obtain the list of tags of a Cluster
//
//	@Summary		Obtain Cluster Tags
//	@Description	Returns a list of Tags belonging to a Cluster given by ID
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	TagListResponse
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id}/tags [get]
func (a APIServer) HandlerGetClusterTags(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Retrieving Cluster's Tags", zap.String("cluster_id", clusterID))

	tags, err := a.sql.GetClusterTags(clusterID)
	if err != nil {
		a.logger.Error("Can't retrieve Tags of cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewTagListResponse(tags))
}

// HandlerPostCluster handles the request for writing a new Cluster in the inventory
//
//	@Summary		Creates a new Cluster in the inventory
//	@Description	Receives and write into the DB the information for a new Cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster	body		inventory.Cluster	true	"New Cluster to be added"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	GenericErrorResponse
//	@Failure		500		{object}	GenericErrorResponse
//	@Router			/clusters [post]
func (a APIServer) HandlerPostCluster(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Capturing request body
	var clusters dto.ClusterDTORequestList
	if err = json.Unmarshal(body, &clusters); err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
		return
	}

	a.logger.Debug("Writing new Clusters", zap.Reflect("clusters", clusters))

	// Filling Account internal ID for every cluster
	var toWriteClusters []inventory.Cluster
	for _, cluster := range clusters.Clusters {
		newCluster := *cluster.ToInventoryCluster()
		if id, err := a.sql.GetAccountInternalID(cluster.AccountID); err != nil {
			a.logger.Error("Can't obtain internal ID for Account", zap.Error(err))
			c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
			return
		} else {
			newCluster.AccountID = id
			toWriteClusters = append(toWriteClusters, newCluster)
		}
	}

	// Writing to DB
	if err = a.sql.WriteClusters(toWriteClusters); err != nil {
		a.logger.Error("Can't write new Clusters into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, responsetypes.PostResponse{
		Count:  len(clusters.Clusters),
		Status: "Cluster(s) Post OK",
	})
}

// HandlerDeleteCluster handles the request for removing a Cluster in the inventory
//
//	@Summary		Deletes a Cluster in the inventory
//	@Description	Deletes a Cluster present in the inventory by its Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id} [delete]
func (a APIServer) HandlerDeleteCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Removing a Cluster", zap.String("cluster_id", clusterID))

	if err := a.sql.DeleteCluster(clusterID); err != nil {
		a.logger.Error("Can't delete Cluster from DB", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, responsetypes.DeleteResponse{
		Count:  1,
		Status: fmt.Sprintf("Cluster '%s' Delete OK", clusterID),
	})
}

// HandlerPatchCluster handles the request for patching a Cluster in the inventory
//
//	@Summary		Patches a Cluster in the inventory
//	@Description	Receives and patch into the DB the information for an existing Cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string				true	"Cluster ID"
//	@Param			cluster		body		inventory.Cluster	true	"Cluster to be modified"
//	@Success		200			{object}	nil
//	@Failure		501			{object}	nil	"Not Implemented"
//	@Router			/clusters/{cluster_id} [patch]
//
// TODO: NOT IMPLEMENTED
func (a APIServer) HandlerPatchCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Patching a Cluster", zap.String("cluster_id", clusterID))

	c.PureJSON(http.StatusNotImplemented, nil)
}
