package handlers

import (
	"net/http"
	"strconv"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/mappers"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
)

// EventHandler handles HTTP requests for events.
type EventHandler struct {
	service services.EventService
}

func NewEventHandler(service services.EventService) *EventHandler {
	return &EventHandler{service: service}
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

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	events, total, err := h.service.ListSystemEvents(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list system events"))
		return
	}

	eventDTOs := mappers.ToSystemEventDTOs(events)
	response := dto.NewListResponse(eventDTOs, total)
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
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
func (h *EventHandler) ListByCluster(c *gin.Context) {
	clusterID := c.Param("id")
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	filters := map[string]interface{}{"resource_id": clusterID}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  filters,
	}

	events, total, err := h.service.ListClusterEvents(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list cluster events"))
		return
	}

	eventDTOs := mappers.ToClusterEventDTOs(events)
	response := dto.NewListResponse(eventDTOs, total)
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}
