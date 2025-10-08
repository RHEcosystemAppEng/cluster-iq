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

// EventHandler handles HTTP requests for events.
type EventHandler struct {
	service services.EventService
	logger  *zap.Logger
}

func NewEventHandler(service services.EventService, logger *zap.Logger) *EventHandler {
	return &EventHandler{
		service: service,
		logger:  logger,
	}
}

type systemEventFilterParams struct {
	TriggeredBy  string `form:"triggered_by"`
	ActionName   string `form:"action_name"`
	ResourceType string `form:"resource_type"`
	Result       string `form:"result"`
	Severity     string `form:"severity"`
}

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

// ListSystem handles the request to list all system events.
//
//	@Summary		List all system events
//	@Description	Returns a paginated list of system events
//	@Tags			Events
//	@Param			page			query		int		false	"Page number for pagination"	default(1)
//	@Param			page_size		query		int		false	"Number of items per page"		default(10)
//	@Param			triggered_by	query		string	false	"Filter by event trigger"
//	@Param			action_name		query		string	false	"Filter by action name"
//	@Param			resource_type	query		string	false	"Filter by resource type"
//	@Param			result			query		string	false	"Filter by event result"
//	@Param			severity		query		string	false	"Filter by event severity"
//	@Success		200				{object}	dto.ListResponse[dto.SystemEvent]
//	@Failure		500				{object}	dto.GenericErrorResponse
//	@Router			/events/system [get]
func (h *EventHandler) ListSystem(c *gin.Context) {
	var req listSystemEventsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
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
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list system events: "+err.Error()))
		return
	}

	response := dto.NewListResponse(db.ToSystemEventDTOResponseList(events), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Create handles the creation of new events.
//
//	@Summary		Create events
//	@Description	Creates one or more new events.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			events	body		[]dto.Newevent	true	"event or events to create"
//	@Success		201			{object}	[]dto.event
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/events [post]
func (h *EventHandler) Create(c *gin.Context) {
	var newEventsDTO dto.EventDTORequest
	if err := c.ShouldBindJSON(&newEventsDTO); err != nil {
		h.logger.Error("error processing received events", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	if _, err := h.service.Create(c.Request.Context(), *newEventsDTO.ToModelEvent()); err != nil {
		h.logger.Error("error creating accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create accounts: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  1,
		Status: "OK"},
	)
}

// ListByCluster handles the request to list all events for a specific cluster.
//
//	@Summary		List all cluster events
//	@Description	Returns a paginated list of cluster events
//	@Tags			Clusters
//	@Param			id			path		string	true	"Cluster ID"
//	@Param			page		query		int		false	"Page number for pagination"	default(1)
//	@Param			page_size	query		int		false	"Number of items per page"		default(10)
//	@Success		200			{object}	dto.ListResponse[dto.ClusterEvent]
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/clusters/{id}/events [get]
func (h *EventHandler) ListCluster(c *gin.Context) {
	clusterID := c.Param("id")
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
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
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list cluster events"))
		return
	}

	response := dto.NewListResponse(db.ToClusterEventDTOResponseList(events), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

func (h *EventHandler) Update(c *gin.Context) {
	var newEventsDTO dto.EventDTORequest
	if err := c.ShouldBindJSON(&newEventsDTO); err != nil {
		h.logger.Error("error processing received events", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	event := newEventsDTO.ToModelEvent()
	if err := h.service.Update(c.Request.Context(), event.ID, event.Result); err != nil {
		h.logger.Error("error updating events", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create accounts: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, responsetypes.PostResponse{
		Count:  1,
		Status: "OK"},
	)
}
