package handlers

import (
	"net/http"
	"strconv"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EventHandler wires HTTP endpoints to the EventService.
type EventHandler struct {
	service services.EventService
	logger  *zap.Logger
}

// NewEventHandler returns an EventHandler with its dependencies.
func NewEventHandler(service services.EventService, logger *zap.Logger) *EventHandler {
	return &EventHandler{
		service: service,
		logger:  logger,
	}
}

// systemEventFilterParams defines the supported filter parameters
type systemEventFilterParams struct {
	TriggeredBy  string `form:"triggered_by"`
	ActionName   string `form:"action_name"`
	ResourceType string `form:"resource_type"`
	Result       string `form:"result"`
	Severity     string `form:"severity"`
}

// toRepoFilters maps bound query params to repository filters.
func (f *systemEventFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.TriggeredBy != "" {
		filters["triggered_by"] = f.TriggeredBy
	}
	if f.ActionName != "" {
		filters["action_name"] = f.ActionName
	}
	if f.ResourceType != "" {
		filters["resource_type"] = f.ResourceType
	}
	if f.Result != "" {
		filters["result"] = f.Result
	}
	if f.Severity != "" {
		filters["severity"] = f.Severity
	}
	return filters
}

type listSystemEventsRequest struct {
	dto.PaginationRequest
	Filters systemEventFilterParams `form:"inline"`
}

// ListSystem returns a paginated list of system events.
//
//	@Summary		List system events
//	@Description	Paginated retrieval with optional filters.
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int		false	"Page number"		default(1)
//	@Param			page_size		query		int		false	"Items per page"	default(10)
//	@Param			triggered_by	query		string	false	"Triggered by"
//	@Param			action_name		query		string	false	"Action name"
//	@Param			resource_type	query		string	false	"Resource type"
//	@Param			result			query		string	false	"Result"
//	@Param			severity		query		string	false	"Severity"
//	@Success		200				{object}	responsetypes.ListResponse[dto.SystemEventDTOResponse]
//	@Failure		400				{object}	responsetypes.GenericErrorResponse
//	@Failure		500				{object}	responsetypes.GenericErrorResponse
//	@Router			/events [get]
func (h *EventHandler) ListSystem(c *gin.Context) {
	var req listSystemEventsRequest

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

	events, total, err := h.service.ListSystemEvents(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing system events", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to list system events: " + err.Error(),
		})
		return
	}

	response := responsetypes.NewListResponse(db.ToSystemEventDTOResponseList(events), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Create creates a new event.
//
//	@Summary		Create event
//	@Description	Create a single event from the request body.
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		dto.EventDTORequest	true	"Event to create"
//	@Success		201		{object}	responsetypes.PostResponse
//	@Failure		400		{object}	responsetypes.GenericErrorResponse
//	@Failure		500		{object}	responsetypes.GenericErrorResponse
//	@Router			/events [post]
//
// NOTE: The handler consumes a single dto.EventDTORequest (not an array).
func (h *EventHandler) Create(c *gin.Context) {
	var newEventsDTO dto.EventDTORequest

	if err := c.ShouldBindJSON(&newEventsDTO); err != nil {
		h.logger.Error("error processing received events", zap.Error(err))
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if _, err := h.service.Create(c.Request.Context(), *newEventsDTO.ToModelEvent()); err != nil {
		h.logger.Error("error creating accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create accounts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  1,
		Status: "OK"},
	)
}

// ListCluster returns a paginated list of events for a cluster.
//
//	@Summary		List cluster events
//	@Description	Paginated events for the specified cluster.
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string	true	"Cluster ID"
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Items per page"	default(10)
//	@Success		200			{object}	responsetypes.ListResponse[dto.ClusterEventDTOResponse]
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/clusters/{id}/events [get]
func (h *EventHandler) ListCluster(c *gin.Context) {
	clusterID := c.Param("id")
	var req dto.PaginationRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid query parameters: " + err.Error(),
		})
		return
	}

	filters := map[string]interface{}{"resource_id": clusterID}

	opts := models.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  filters,
	}

	events, total, err := h.service.ListClusterEvents(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error getting events for cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to list cluster events",
		})
		return
	}

	response := responsetypes.NewListResponse(db.ToClusterEventDTOResponseList(events), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Update modifies an event's result based on the payload.
//
//	@Summary		Update event
//	@Description	Update an event result using the provided payload.
//	@Tags			Events
//	@Accept			json
//	@Produce		json
//	@Param			event	body		dto.EventDTORequest	true	"Event update payload"
//	@Success		200		{object}	responsetypes.PostResponse
//	@Failure		400		{object}	responsetypes.GenericErrorResponse
//	@Failure		500		{object}	responsetypes.GenericErrorResponse
//	@Router			/events [patch]
func (h *EventHandler) Update(c *gin.Context) {
	var newEventsDTO dto.EventDTORequest

	if err := c.ShouldBindJSON(&newEventsDTO); err != nil {
		h.logger.Error("error processing received events", zap.Error(err))
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	event := newEventsDTO.ToModelEvent()
	if err := h.service.Update(c.Request.Context(), event.ID, event.Result); err != nil {
		h.logger.Error("error updating events", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create accounts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responsetypes.PostResponse{
		Count:  1,
		Status: "OK"},
	)
}
