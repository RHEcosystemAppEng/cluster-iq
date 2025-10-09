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

// ActionHandler wires HTTP endpoints to the ActionService.
type ActionHandler struct {
	service services.ActionService
	logger  *zap.Logger
}

// NewActionHandler returns an ActionHandler with its dependencies.
func NewActionHandler(service services.ActionService, logger *zap.Logger) *ActionHandler {
	return &ActionHandler{
		service: service,
		logger:  logger,
	}
}

// scheduledActionFilterParams defines the supported filter parameters
type scheduledActionFilterParams struct {
	Enabled string `form:"enabled"`
	Status  string `form:"status"`
}

// toRepoFilters maps bound query parameters to repository filters.
func (f *scheduledActionFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.Enabled != "" {
		filters["enabled"] = f.Enabled
	}
	if f.Status != "" {
		filters["status"] = f.Status
	}
	return filters
}

// listScheduledActionsRequest carries pagination and filters for List.
type listScheduledActionsRequest struct {
	dto.PaginationRequest
	Filters scheduledActionFilterParams `form:"inline"`
}

// List returns a paginated list of scheduled actions.
//
//	@Summary		List scheduled actions
//	@Description	Paginated retrieval with optional filters.
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			enabled		query		string	false	"Enabled state filter (true/false)"
//	@Param			status		query		string	false	"Status filter"
//	@Param			page		query		int		false	"Page number"			default(1)
//	@Param			page_size	query		int		false	"Items per page"		default(10)
//	@Success		200			{object}	dto.ListResponse[dto.ScheduledAction]
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/schedule [get]
func (h *ActionHandler) List(c *gin.Context) {
	var req listScheduledActionsRequest

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

	actions, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing actions schedule", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to list scheduled actions",
		})
		return
	}

	response := responsetypes.NewListResponse(db.ToActionDTOResponseList(actions), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Get returns a scheduled action by ID.
//
//	@Summary		Get scheduled action by ID
//	@Description	Return a scheduled action resource.
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Scheduled action ID"
//	@Success		200	{object}	dto.ScheduledAction
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/schedule/{id} [get]
func (h *ActionHandler) Get(c *gin.Context) {
	actionID := c.Param("id")

	action, err := h.service.Get(c.Request.Context(), actionID)
	if err != nil {
		h.logger.Error("error getting action", zap.String("action_id", actionID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Scheduled action not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve scheduled action",
		})
		return
	}

	c.JSON(http.StatusOK, action.ToActionDTOResponse())
}

// Create creates one or more actions.
//
//	@Summary		Create actions
//	@Description	Create one or multiple actions.
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			actions	body		[]dto.ActionDTORequest	true	"Actions to create"
//	@Success		201			{object}	responsetypes.PostResponse
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/actions [post]
func (h *ActionHandler) Create(c *gin.Context) {
	var newActionsDTO []dto.ActionDTORequest

	if err := c.ShouldBindJSON(&newActionsDTO); err != nil {
		h.logger.Error("error processing received actions", zap.Error(err))
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	actions := dto.ToModelActionList(newActionsDTO)
	if actions == nil {
		err := errors.New("error when processing actions")
		h.logger.Error("error processing actions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create actions: " + err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), *actions); err != nil {
		h.logger.Error("error creating actions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create actions: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(newActionsDTO),
		Status: "OK"},
	)
}

// Enable activates a scheduled action by ID.
//
//	@Summary		Enable scheduled action
//	@Description	Enable a scheduled action.
//	@Tags			Actions
//	@Param			id	path		string	true	"Scheduled action ID"
//	@Success		200	{object}	nil
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/schedule/{id}/enable [patch]
func (h *ActionHandler) Enable(c *gin.Context) {
	actionID := c.Param("id")

	if err := h.service.Enable(c.Request.Context(), actionID); err != nil {
		h.logger.Error("error enabling action", zap.String("action_id", actionID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to enable scheduled action",
		})
		return
	}

	c.Status(http.StatusOK)
}

// Disable deactivates a scheduled action by ID.
//
//	@Summary		Disable scheduled action
//	@Description	Disable a scheduled action.
//	@Tags			Actions
//	@Param			id	path		string	true	"Scheduled action ID"
//	@Success		200	{object}	nil
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/schedule/{id}/disable [patch]
func (h *ActionHandler) Disable(c *gin.Context) {
	actionID := c.Param("id")

	if err := h.service.Disable(c.Request.Context(), actionID); err != nil {
		h.logger.Error("error disabling action", zap.String("action_id", actionID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to disable scheduled action",
		})
		return
	}

	c.Status(http.StatusOK)
}

// Delete removes an action by ID.
//
//	@Summary		Delete an action
//	@Description	Delete an action by ID.
//	@Tags			Actions
//	@Accept			json
//	@Param			id	path		string	true	"Action ID"
//	@Success		204		{object}	nil
//	@Failure		404		{object}	responsetypes.GenericErrorResponse
//	@Failure		500		{object}	responsetypes.GenericErrorResponse
//	@Router			/actions/{id} [delete]
func (h *ActionHandler) Delete(c *gin.Context) {
	actionID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), actionID); err != nil {
		h.logger.Error("error deleting action", zap.String("action_id", actionID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "action not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to delete Action: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
