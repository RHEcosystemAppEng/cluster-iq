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

// InstanceHandler handles HTTP requests for instances.
type InstanceHandler struct {
	service services.InstanceService
}

func NewInstanceHandler(service services.InstanceService) *InstanceHandler {
	return &InstanceHandler{service: service}
}

type instanceFilterParams struct {
	ClusterID string `form:"cluster_id"`
	Status    string `form:"status"`
}

func (f *instanceFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.ClusterID != "" {
		filters["cluster_id"] = f.ClusterID
	}
	if f.Status != "" {
		filters["status"] = f.Status
	}
	return filters
}

type listInstancesRequest struct {
	dto.PaginationRequest
	Filters instanceFilterParams `form:"inline"`
}

// List handles the request for obtaining the Instance list.
//
//	@Summary		List instances
//	@Description	Returns a paginated list of instances based on optional filters.
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number for pagination"	default(1)
//	@Param			page_size	query		int		false	"Number of items per page"		default(10)
//	@Param			cluster_id	query		string	false	"Filter by cluster ID"
//	@Param			status		query		string	false	"Filter by instance status"
//	@Success		200			{object}	dto.ListResponse[dto.Instance]
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/instances [get]
func (h *InstanceHandler) List(c *gin.Context) {
	var req listInstancesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	instances, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve instances"))
		return
	}

	instanceDTOs := mappers.ToInstanceDTOs(instances)
	response := dto.NewListResponse(instanceDTOs, total)
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Get handles the request for obtaining a single instance by its ID.
//
//	@Summary		Get an instance by ID
//	@Description	Returns a single instance.
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Instance ID"
//	@Success		200	{object}	dto.Instance
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/instances/{id} [get]
func (h *InstanceHandler) Get(c *gin.Context) {
	instanceID := c.Param("id")

	instance, err := h.service.Get(c.Request.Context(), instanceID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Instance not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve instance"))
		return
	}

	instanceDTO := mappers.ToInstanceDTO(instance)
	c.JSON(http.StatusOK, instanceDTO)
}

// GetSummary handles the request for obtaining a summary of instance statuses.
//
//	@Summary		Get instances summary
//	@Description	Returns a summary of instance counts.
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	inventory.InstancesSummary
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/instances/summary [get]
// func (h *InstanceHandler) GetSummary(c *gin.Context) {
// 	summary, err := h.service.GetSummary(c.Request.Context())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve instances summary"))
// 		return
// 	}

// 	c.JSON(http.StatusOK, summary)
// }

// ListByCluster handles the request for obtaining instances for a specific cluster.
//
//	@Summary		List cluster's instances
//	@Description	Returns a paginated list of instances for a specific cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Cluster ID"
//	@Param			page		query		int		false	"Page number for pagination"	default(1)
//	@Param			page_size	query		int		false	"Number of items per page"		default(10)
//	@Success		200			{object}	dto.ListResponse[dto.Instance]
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id}/instances [get]
func (h *InstanceHandler) ListByCluster(c *gin.Context) {
	clusterID := c.Param("id")
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	filters := map[string]interface{}{"cluster_id": clusterID}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  filters,
	}

	instances, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve instances for cluster "+clusterID))
		return
	}

	instanceDTOs := mappers.ToInstanceDTOs(instances)
	response := dto.NewListResponse(instanceDTOs, total)
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}
