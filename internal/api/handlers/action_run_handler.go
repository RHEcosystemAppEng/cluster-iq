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

// ActionRunHandler wires HTTP endpoints to the ActionRunService.
type ActionRunHandler struct {
	service services.ActionRunService
	logger  *zap.Logger
}

// NewActionRunHandler returns an ActionRunHandler with its dependencies.
func NewActionRunHandler(service services.ActionRunService, logger *zap.Logger) *ActionRunHandler {
	return &ActionRunHandler{
		service: service,
		logger:  logger,
	}
}

type actionRunFilterParams struct {
	ScheduleID string `form:"schedule_id"`
	Status     string `form:"status"`
}

func (f *actionRunFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.ScheduleID != "" {
		filters["schedule_id"] = f.ScheduleID
	}
	if f.Status != "" {
		filters["status"] = f.Status
	}
	return filters
}

type listActionRunsRequest struct {
	dto.PaginationRequest
	Filters actionRunFilterParams `form:"inline"`
}

// List returns a paginated list of action runs.
//
//	@Summary		List action runs
//	@Description	Paginated retrieval of action execution history.
//	@Tags			ActionRuns
//	@Accept			json
//	@Produce		json
//	@Param			schedule_id	query		string	false	"Filter by schedule ID"
//	@Param			status		query		string	false	"Filter by status"
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Items per page"	default(10)
//	@Success		200			{object}	dto.ActionRunListResponse
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/action-runs [get]
func (h *ActionRunHandler) List(c *gin.Context) {
	var req listActionRunsRequest

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

	runs, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing action runs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to list action runs",
		})
		return
	}

	response := responsetypes.NewListResponse((&convert.ConverterImpl{}).ToActionRunDTOs(runs), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Get returns an action run by ID.
//
//	@Summary		Get action run by ID
//	@Description	Return an action run record.
//	@Tags			ActionRuns
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Action run ID"
//	@Success		200	{object}	dto.ActionRunDTOResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/action-runs/{id} [get]
func (h *ActionRunHandler) Get(c *gin.Context) {
	runID := c.Param("id")

	run, err := h.service.Get(c.Request.Context(), runID)
	if err != nil {
		h.logger.Error("error getting action run", zap.String("run_id", runID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Action run not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve action run",
		})
		return
	}

	c.JSON(http.StatusOK, (&convert.ConverterImpl{}).ToActionRunDTO(run))
}

// Create creates a new action run.
//
//	@Summary		Create action run
//	@Description	Create a new execution record for a scheduled action.
//	@Tags			ActionRuns
//	@Accept			json
//	@Produce		json
//	@Param			run	body		dto.ActionRunDTORequest	true	"Action run to create"
//	@Success		201	{object}	responsetypes.PostResponse
//	@Failure		400	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/action-runs [post]
func (h *ActionRunHandler) Create(c *gin.Context) {
	var req dto.ActionRunDTORequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	id, err := h.service.Create(c.Request.Context(), req.ScheduleID)
	if err != nil {
		h.logger.Error("error creating action run", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create action run: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  1,
		Status: strconv.FormatInt(id, 10),
	})
}

// Update updates an action run's status.
//
//	@Summary		Update action run
//	@Description	Update the status of an action run (called by scanner on completion).
//	@Tags			ActionRuns
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"Action run ID"
//	@Param			run	body		dto.ActionRunDTORequest	true	"Updated status"
//	@Success		200	{object}	nil
//	@Failure		400	{object}	responsetypes.GenericErrorResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/action-runs/{id} [patch]
func (h *ActionRunHandler) Update(c *gin.Context) {
	runID := c.Param("id")

	var req dto.ActionRunDTORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if _, err := h.service.Get(c.Request.Context(), runID); err != nil {
		h.logger.Error("error updating action run", zap.String("run_id", runID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Action run not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to update action run",
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), runID, req.Status, req.ErrorMsg); err != nil {
		h.logger.Error("error updating action run", zap.String("run_id", runID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to update action run: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
