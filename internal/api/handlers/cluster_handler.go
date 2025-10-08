package handlers

import (
	"errors"
	"net/http"
	"strconv"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ClusterHandler handles HTTP requests for clusters.
type ClusterHandler struct {
	service services.ClusterService
	logger  *zap.Logger
}

func NewClusterHandler(service services.ClusterService, logger *zap.Logger) *ClusterHandler {
	return &ClusterHandler{
		service: service,
		logger:  logger,
	}
}

type clusterFilterParams struct {
	Status   string `form:"status"`
	Provider string `form:"provider"`
	Region   string `form:"region"`
	Account  string `form:"account"`
}

func (f *clusterFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.Status != "" {
		filters["status"] = f.Status
	}
	if f.Provider != "" {
		filters["provider"] = f.Provider
	}
	if f.Region != "" {
		filters["region"] = f.Region
	}
	if f.Account != "" {
		filters["account_name"] = f.Account
	}
	return filters
}

type listClustersRequest struct {
	dto.PaginationRequest
	Filters clusterFilterParams `form:"inline"`
}

// List handles the request for obtaining the Cluster list.
//
//	@Summary		List clusters
//	@Description	Returns a paginated list of clusters based on optional filters.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number for pagination"		default(1)
//	@Param			page_size	query		int		false	"Number of items per page"			default(10)
//	@Param			status		query		string	false	"Filter by cluster status"			example(Running)
//	@Param			provider	query		string	false	"Filter by cloud provider"			example(aws)
//	@Param			region		query		string	false	"Filter by cloud provider region"	example(us-east-1)
//	@Param			account		query		string	false	"Filter by account name"
//	@Success		200			{object}	dto.ListResponse[dto.Cluster]
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/clusters [get]
func (h *ClusterHandler) List(c *gin.Context) {
	var req listClustersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	opts := models.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	clusters, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing clusters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve clusters"))
		return
	}

	response := dto.NewListResponse(db.ToClusterDTOResponseList(clusters), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Get handles the request for obtaining a single cluster by its ID.
//
//	@Summary		Get a cluster by ID
//	@Description	Returns a single cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		200	{object}	dto.Cluster
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id} [get]
func (h *ClusterHandler) Get(c *gin.Context) {
	clusterID := c.Param("id")

	cluster, err := h.service.Get(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error getting a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve cluster"))
		return
	}

	c.JSON(http.StatusOK, cluster.ToClusterDTOResponse())
}

// GetInstances handles the request for obtaining a single cluster by its ID and its instances.
//
//	@Summary		Get a cluster by ID with its instances
//	@Description	Returns a single cluster with its instances.
//	@Tags			Clusters Instances
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		200	{object}	dto.Cluster
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id} [get]
func (h *ClusterHandler) GetInstances(c *gin.Context) {
	clusterID := c.Param("id")

	instances, err := h.service.GetInstances(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error getting a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve cluster"))
		return
	}

	c.JSON(http.StatusOK, db.ToInstanceDTOResponseList(instances))
}

// Create handles the creation of new clusters.
//
//	@Summary		Create clusters
//	@Description	Creates one or more new clusters.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			clusters	body		[]dto.Cluster	true	"A list of clusters to create"
//	@Success		201			{object}	nil
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/clusters [post]
func (h *ClusterHandler) Create(c *gin.Context) {
	var newClusterDTOs []dto.ClusterDTORequest
	if err := c.ShouldBindJSON(&newClusterDTOs); err != nil {
		h.logger.Error("error processing cluster creation request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	if err := h.service.Create(c.Request.Context(), *dto.ToInventoryClusterList(newClusterDTOs)); err != nil {
		h.logger.Error("error creating a cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create clusters: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(newClusterDTOs),
		Status: "OK"},
	)
}

// Delete handles the deletion of a cluster.
//
//	@Summary		Delete a cluster
//	@Description	Deletes a cluster by its ID.
//	@Tags			Clusters
//	@Accept			json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id} [delete]
func (h *ClusterHandler) Delete(c *gin.Context) {
	clusterID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), clusterID); err != nil {
		h.logger.Error("error deleting a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to delete cluster: "+err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// PowerOn handles the request to power on a cluster.
//
//	@Summary		Power on a cluster
//	@Description	Sends a request to power on a cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		202	{object}	dto.GenericResponse
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id}/power_on [post]
func (h *ClusterHandler) PowerOn(c *gin.Context) {
	clusterID := c.Param("id")
	err := h.service.PowerOn(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error powering on a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to power on cluster: "+err.Error()))
		return
	}

	c.JSON(http.StatusAccepted, dto.NewGenericResponse("Power on request accepted"))
}

// PowerOff handles the request to power off a cluster.
//
//	@Summary		Power off a cluster
//	@Description	Sends a request to power off a cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		202	{object}	dto.GenericResponse
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id}/power_off [post]
func (h *ClusterHandler) PowerOff(c *gin.Context) {
	clusterID := c.Param("id")
	err := h.service.PowerOff(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error powering off a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to power off cluster: "+err.Error()))
		return
	}

	c.JSON(http.StatusAccepted, dto.NewGenericResponse("Power off request accepted"))
}

// GetTags handles the request to retrieve tags for a specific cluster.
//
//	@Summary		Get cluster tags
//	@Description	Returns all tags for a specific cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		200	{object}	[]dto.Tag
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id}/tags [get]
func (h *ClusterHandler) GetTags(c *gin.Context) {
	clusterID := c.Param("id")
	tags, err := h.service.GetTags(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error getting cluster tags", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to get cluster tags: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, db.ToTagsDTOResponseList(tags))
}

// Update handles the request to update an existing cluster.
//
//	@Summary		Update a cluster
//	@Description	Updates mutable fields of an existing cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string		true	"Cluster ID"
//	@Param			cluster	body		dto.Cluster	true	"Cluster data to update"
//	@Success		200		{object}	nil
//	@Failure		501		{object}	nil	"Not Implemented"
//	@Router			/clusters/{id} [patch]
func (h *ClusterHandler) Update(c *gin.Context) {
	// TODO
	c.PureJSON(http.StatusNotImplemented, nil)
}
