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

// InstanceHandler handles HTTP requests for instances.
type InstanceHandler struct {
	service services.InstanceService
	logger  *zap.Logger
}

func NewInstanceHandler(service services.InstanceService, logger *zap.Logger) *InstanceHandler {
	return &InstanceHandler{
		service: service,
		logger:  logger,
	}
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

	opts := models.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	instances, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve instances"))
		return
	}

	response := dto.NewListResponse(db.ToInstanceDTOResponseList(instances), total)

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
		h.logger.Error("error getting instance", zap.String("instance_id", instanceID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Instance not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve instance"))
		return
	}

	c.JSON(http.StatusOK, instance.ToInstanceDTOResponse())
}

// Create handles the creation of new instances.
//
//	@Summary		Create instances
//	@Description	Creates one or more new instances.
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			instances	body		[]dto.Instance	true	"A list of instances to create"
//	@Success		201			{object}	nil
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/instances [post]
func (h *InstanceHandler) Create(c *gin.Context) {
	var newInstanceDTOs []dto.InstanceDTORequest
	if err := c.ShouldBindJSON(&newInstanceDTOs); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	if err := h.service.Create(c.Request.Context(), *dto.ToInventoryInstanceList(newInstanceDTOs)); err != nil {
		h.logger.Error("error creating instances", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create instances: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(newInstanceDTOs),
		Status: "OK"},
	)
}
