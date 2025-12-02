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

// InstanceHandler exposes instance-related HTTP endpoints.
type InstanceHandler struct {
	service services.InstanceService
	logger  *zap.Logger
}

// NewInstanceHandler wires service and logger into the handler.
func NewInstanceHandler(service services.InstanceService, logger *zap.Logger) *InstanceHandler {
	return &InstanceHandler{
		service: service,
		logger:  logger,
	}
}

// instanceFilterParams defines the supported filter parameters
type instanceFilterParams struct {
	ClusterID string `form:"cluster_id"`
	Status    string `form:"status"`
}

// toRepoFilters translates bound query params into repository filters.
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

// listInstancesRequest carries pagination and filters for List.
type listInstancesRequest struct {
	dto.PaginationRequest
	Filters instanceFilterParams `form:"inline"`
}

// List returns a paginated list of instances.
//
//	@Summary		List instances
//	@Description	Paginated retrieval with optional filters.
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Items per page"	default(10)
//	@Param			cluster_id	query		string	false	"Cluster ID filter"
//	@Param			status		query		string	false	"Instance status filter"
//	@Success		200			{object}	responsetypes.ListResponse[dto.InstanceDTOResponse]
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/instances [get]
func (h *InstanceHandler) List(c *gin.Context) {
	var req listInstancesRequest

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

	instances, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve instances",
		})
		return
	}

	response := responsetypes.NewListResponse(db.ToInstanceDTOResponseList(instances), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Get returns a single instance by ID.
//
//	@Summary		Get instance by ID
//	@Description	Return a single instance resource.
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Instance ID"
//	@Success		200	{object}	dto.InstanceDTOResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/instances/{id} [get]
func (h *InstanceHandler) Get(c *gin.Context) {
	instanceID := c.Param("id")

	instance, err := h.service.Get(c.Request.Context(), instanceID)
	if err != nil {
		h.logger.Error("error getting instance", zap.String("instance_id", instanceID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Instance not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve instance",
		})
		return
	}

	c.JSON(http.StatusOK, instance.ToInstanceDTOResponse())
}

// Create inserts one or more instances.
//
//	@Summary		Create instances
//	@Description	Create one or multiple instances.
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			instances	body		[]dto.InstanceDTORequest	true	"Instances to create"
//	@Success		201			{object}	responsetypes.PostResponse
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/instances [post]
func (h *InstanceHandler) Create(c *gin.Context) {
	var newInstanceDTOs []dto.InstanceDTORequest

	if err := c.ShouldBindJSON(&newInstanceDTOs); err != nil {
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), *dto.ToInventoryInstanceList(newInstanceDTOs)); err != nil {
		h.logger.Error("error creating instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create instances: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(newInstanceDTOs),
		Status: "OK"},
	)
}
