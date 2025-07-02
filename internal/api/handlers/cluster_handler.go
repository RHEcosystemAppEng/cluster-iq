package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/mappers"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
)

// ClusterHandler handles HTTP requests for clusters.
type ClusterHandler struct {
	service services.ClusterService
}

func NewClusterHandler(service services.ClusterService) *ClusterHandler {
	return &ClusterHandler{service: service}
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
//	@Param			page		query	int		false	"Page number for pagination"			default(1)
//	@Param			pageSize	query	int		false	"Number of items per page"				default(10)
//	@Param			status		query	string	false	"Filter by cluster status"				example(Running)
//	@Param			provider	query	string	false	"Filter by cloud provider"				example(aws)
//	@Param			region		query	string	false	"Filter by cloud provider region"		example(us-east-1)
//	@Param			account		query	string	false	"Filter by account name"
//	@Success		200			{object}	dto.ListResponse[dto.Cluster]
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/clusters [get]
func (h *ClusterHandler) List(c *gin.Context) {
	var req listClustersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	clusters, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve clusters"))
		return
	}

	clusterDTOs := mappers.ToClusterDTOs(clusters)
	response := dto.NewListResponse(clusterDTOs, total)

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
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/clusters/{id} [get]
func (h *ClusterHandler) Get(c *gin.Context) {
	clusterID := c.Param("id")

	cluster, err := h.service.Get(c.Request.Context(), clusterID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve cluster"))
		return
	}

	clusterDTO := mappers.ToClusterDTO(cluster)
	c.JSON(http.StatusOK, clusterDTO)
}

// GetSummary handles the request for obtaining a summary of cluster statuses.
//
//	@Summary		Get cluster summary
//	@Description	Returns a summary of cluster counts by status.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	inventory.ClustersSummary
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/clusters/summary [get]
func (h *ClusterHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve cluster summary"))
		return
	}

	c.JSON(http.StatusOK, summary)
}

// Create handles the creation of a new cluster.
//
//	@Summary		Create a cluster
//	@Description	Creates a new cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster	body		dto.Cluster	true	"Cluster to create"
//	@Success		201		{object}	dto.Cluster
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/clusters [post]
func (h *ClusterHandler) Create(c *gin.Context) {
	var newClusterDTO dto.Cluster
	if err := c.ShouldBindJSON(&newClusterDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	cluster := mappers.ToClusterModel(newClusterDTO)
	err := h.service.Create(c.Request.Context(), cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create cluster: "+err.Error()))
		return
	}

	createdClusterDTO := mappers.ToClusterDTO(cluster)
	c.JSON(http.StatusCreated, createdClusterDTO)
}

// Delete handles the deletion of a cluster.
//
//	@Summary		Delete a cluster
//	@Description	Deletes a cluster by its ID.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Cluster ID"
//	@Success		204	{object}	nil
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/clusters/{id} [delete]
func (h *ClusterHandler) Delete(c *gin.Context) {
	clusterID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), clusterID); err != nil {
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
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/clusters/{id}/power_on [post]
func (h *ClusterHandler) PowerOn(c *gin.Context) {
	clusterID := c.Param("id")
	err := h.service.PowerOn(c.Request.Context(), clusterID)

	if err != nil {
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
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/clusters/{id}/power_off [post]
func (h *ClusterHandler) PowerOff(c *gin.Context) {
	clusterID := c.Param("id")
	err := h.service.PowerOff(c.Request.Context(), clusterID)

	if err != nil {
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
//	@Failure		404	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/clusters/{id}/tags [get]
func (h *ClusterHandler) GetTags(c *gin.Context) {
	clusterID := c.Param("id")
	tags, err := h.service.GetTags(c.Request.Context(), clusterID)

	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to get cluster tags: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, mappers.ToTagDTOs(tags))
}

// Update handles the request to update an existing cluster.
//
//	@Summary		Update a cluster
//	@Description	Updates mutable fields of an existing cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string		true	"Cluster ID"
//	@Param			cluster	body	dto.Cluster	true	"Cluster data to update"
//	@Success		200		{object}	dto.Cluster
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/clusters/{id} [patch]
func (h *ClusterHandler) Update(c *gin.Context) {
	clusterID := c.Param("id")

	var clusterDTO dto.Cluster
	if err := c.ShouldBindJSON(&clusterDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	// Ensure the ID from the path is used, not from the body
	clusterDTO.ID = clusterID

	// Fetch the existing cluster to ensure it exists
	existingCluster, err := h.service.Get(c.Request.Context(), clusterID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Cluster not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve cluster for update: "+err.Error()))
		return
	}

	// Update only the mutable fields from the DTO
	// Here we assume ToClusterModel handles partial updates gracefully or we do it manually
	updatedClusterModel := mappers.ToClusterModel(clusterDTO)
	// Preserve non-mutable fields from the existing model
	updatedClusterModel.Provider = existingCluster.Provider
	updatedClusterModel.Status = existingCluster.Status

	if err := h.service.Update(c.Request.Context(), updatedClusterModel); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to update cluster: "+err.Error()))
		return
	}

	// Return the fully updated object
	finalCluster, err := h.service.Get(c.Request.Context(), clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve updated cluster: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, mappers.ToClusterDTO(finalCluster))
}
