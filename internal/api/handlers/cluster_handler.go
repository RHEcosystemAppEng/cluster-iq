package handlers

import (
	"errors"
	"net/http"
	"strconv"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/convert"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ClusterHandler connects HTTP endpoints to the ClusterService.
type ClusterHandler struct {
	service services.ClusterService
	logger  *zap.Logger
}

// NewClusterHandler returns a ClusterHandler with its dependencies.
func NewClusterHandler(service services.ClusterService, logger *zap.Logger) *ClusterHandler {
	return &ClusterHandler{
		service: service,
		logger:  logger,
	}
}

// accountFilterParams defines the supported filter parameters
type clusterFilterParams struct {
	Status   string `form:"status"`
	Provider string `form:"provider"`
	Region   string `form:"region"`
	Account  string `form:"account"`
}

// toRepoFilters maps bound query params to repository filters.
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

// listClustersRequest encapsulates pagination and filters for List.
type listClustersRequest struct {
	dto.PaginationRequest
	Filters clusterFilterParams `form:"inline"`
}

// List returns a paginated list of clusters.
//
//	@Summary		List clusters
//	@Description	Paginated retrieval with optional filters.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Items per page"	default(10)
//	@Param			status		query		string	false	"Cluster status"
//	@Param			provider	query		string	false	"Cloud provider"
//	@Param			region		query		string	false	"Provider region"
//	@Param			account		query		string	false	"Account name"
//	@Success		200			{object}	dto.ClusterListResponse
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters [get]
func (h *ClusterHandler) List(c *gin.Context) {
	var req listClustersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid query parameters: " + err.Error(),
		})
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
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve clusters",
		})
		return
	}

	response := responsetypes.NewListResponse((&convert.ConverterImpl{}).ToClusterDTOs(clusters), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Get returns a cluster by ID.
//
//	@Summary		Get cluster by ID
//	@Description	Return a single cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		200	{object}	dto.ClusterDTOResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters/{id} [get]
func (h *ClusterHandler) Get(c *gin.Context) {
	clusterID := c.Param("id")

	cluster, err := h.service.Get(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error getting a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Cluster not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve cluster",
		})
		return
	}

	c.JSON(http.StatusOK, (&convert.ConverterImpl{}).ToClusterDTO(*cluster))
}

// GetInstances returns the instances for a cluster ID.
//
//	@Summary		Get cluster instances
//	@Description	Return instances belonging to the specified cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		200	{object}	[]dto.InstanceDTOResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters/{id}/instances [get]
func (h *ClusterHandler) GetInstances(c *gin.Context) {
	clusterID := c.Param("id")

	instances, err := h.service.GetInstances(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error getting a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Cluster not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve cluster",
		})
		return
	}

	c.JSON(http.StatusOK, (&convert.ConverterImpl{}).ToInstanceDTOs(instances))
}

// Create creates one or more clusters.
//
//	@Summary		Create clusters
//	@Description	Create one or multiple clusters.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			clusters	body		[]dto.ClusterDTORequest	true	"Clusters to create"
//	@Success		201			{object}	responsetypes.PostResponse
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters [post]
func (h *ClusterHandler) Create(c *gin.Context) {
	var newClusterDTOs []dto.ClusterDTORequest

	if err := c.ShouldBindJSON(&newClusterDTOs); err != nil {
		h.logger.Error("error processing cluster creation request", zap.Error(err))
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), *dto.ToInventoryClusterList(newClusterDTOs)); err != nil {
		h.logger.Error("error creating a cluster", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create clusters: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(newClusterDTOs),
		Status: "OK",
	},
	)
}

// Delete removes a cluster by ID.
//
//	@Summary		Delete a cluster
//	@Description	Delete a cluster by ID.
//	@Tags			Clusters
//	@Accept			json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters/{id} [delete]
func (h *ClusterHandler) Delete(c *gin.Context) {
	clusterID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), clusterID); err != nil {
		h.logger.Error("error deleting a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Cluster not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to delete cluster: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// PowerOn triggers a power-on operation for a cluster.
//
//	@Summary		Power on a cluster
//	@Description	Request powering on a cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		202	{object}	responsetypes.GenericResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters/{id}/power_on [post]
func (h *ClusterHandler) PowerOn(c *gin.Context) {
	clusterID := c.Param("id")

	if err := h.service.PowerOn(c.Request.Context(), clusterID); err != nil {
		h.logger.Error("error powering on a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Cluster not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to power on cluster: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, responsetypes.GenericResponse{
		Message: "Power on request accepted",
	})
}

// PowerOff triggers a power-off operation for a cluster.
//
//	@Summary		Power off a cluster
//	@Description	Request powering off a cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		202	{object}	responsetypes.GenericResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters/{id}/power_off [post]
func (h *ClusterHandler) PowerOff(c *gin.Context) {
	clusterID := c.Param("id")

	if err := h.service.PowerOff(c.Request.Context(), clusterID); err != nil {
		h.logger.Error("error powering off a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Cluster not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to power off cluster: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, responsetypes.GenericResponse{
		Message: "Power off request accepted",
	})
}

// GetTags returns tags for the specified cluster.
//
//	@Summary		Get cluster tags
//	@Description	Return all tags for a cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		200	{object}	[]dto.TagDTOResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters/{id}/tags [get]
func (h *ClusterHandler) GetTags(c *gin.Context) {
	clusterID := c.Param("id")

	tags, err := h.service.GetTags(c.Request.Context(), clusterID)
	if err != nil {
		h.logger.Error("error getting cluster tags", zap.String("cluster_id", clusterID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Cluster not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to get cluster tags: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, (&convert.ConverterImpl{}).ToTagDTOs(tags))
}

// Update applies partial updates to a cluster.
//
//	@Summary		Update a cluster
//	@Description	Patch mutable fields of a cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Cluster ID"
//	@Param			cluster	body		dto.ClusterDTOResponse	true	"Partial cluster payload"
//	@Success		200		{object}	nil
//	@Failure		501		{object}	nil	"Not Implemented"
//	@Router			/clusters/{id} [patch]
func (h *ClusterHandler) Update(c *gin.Context) {
	// TODO: Implement partial update strategy
	c.PureJSON(http.StatusNotImplemented, nil)
}
