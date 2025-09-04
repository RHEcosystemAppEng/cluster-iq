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

// ActionHandler handles HTTP requests for actions.
type ActionHandler struct {
	service services.ActionService
}

// NewActionHandler creates a new ActionHandler.
func NewActionHandler(service services.ActionService) *ActionHandler {
	return &ActionHandler{service: service}
}

type scheduledActionFilterParams struct {
	Enabled string `form:"enabled"`
	Status  string `form:"status"`
}

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

type listScheduledActionsRequest struct {
	dto.PaginationRequest
	Filters scheduledActionFilterParams `form:"inline"`
}

// ListScheduled handles the request to list all scheduled actions.
//
//	@Summary		List all scheduled actions
//	@Description	Returns a list of scheduled actions
//	@Tags			Actions
//	@Param			enabled		query		string	false	"Filter by enabled state (true/false)"
//	@Param			status		query		string	false	"Filter by action status"
//	@Param			page		query		int		false	"Page number for pagination"	default(1)
//	@Param			page_size	query		int		false	"Number of items per page"		default(10)
//	@Success		200			{object}	dto.ListResponse[dto.ScheduledAction]
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/schedule [get]
func (h *ActionHandler) ListScheduled(c *gin.Context) {
	var req listScheduledActionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	actions, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list scheduled actions"))
		return
	}

	actionDTOs := mappers.ToScheduledActionDTOList(actions)
	response := dto.NewListResponse(actionDTOs, total)
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// GetScheduled handles the request to get a single scheduled action by its ID.
//
//	@Summary		Get scheduled action by ID
//	@Description	Returns details of a specific scheduled action identified by its ID.
//	@Tags			Actions
//	@Param			id	path		string	true	"Scheduled action identifier"
//	@Success		200	{object}	dto.ScheduledAction
//	@Failure		404	{object}	dto.GenericErrorResponse
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/schedule/{id} [get]
func (h *ActionHandler) GetScheduled(c *gin.Context) {
	actionID := c.Param("id")

	action, err := h.service.Get(c.Request.Context(), actionID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Scheduled action not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve scheduled action"))
		return
	}

	actionDTO := mappers.ToScheduledActionDTO(action)
	c.JSON(http.StatusOK, actionDTO)
}

// EnableScheduled handles the request to enable a scheduled action.
//
//	@Summary		Enable scheduled action
//	@Description	Activates a scheduled action specified by its ID.
//	@Tags			Actions
//	@Param			id	path		string	true	"Scheduled action identifier"
//	@Success		204	{object}	nil
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/schedule/{id}/enable [patch]
func (h *ActionHandler) EnableScheduled(c *gin.Context) {
	actionID := c.Param("id")

	if err := h.service.Enable(c.Request.Context(), actionID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to enable scheduled action"))
		return
	}

	c.Status(http.StatusNoContent)
}

// DisableScheduled handles the request to disable a scheduled action.
//
//	@Summary		Disable scheduled action
//	@Description	Deactivates a scheduled action specified by its ID.
//	@Tags			Actions
//	@Param			id	path		string	true	"Scheduled action identifier"
//	@Success		204	{object}	nil
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/schedule/{id}/disable [patch]
func (h *ActionHandler) DisableScheduled(c *gin.Context) {
	actionID := c.Param("id")

	if err := h.service.Disable(c.Request.Context(), actionID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to disable scheduled action"))
		return
	}

	c.Status(http.StatusNoContent)
}
