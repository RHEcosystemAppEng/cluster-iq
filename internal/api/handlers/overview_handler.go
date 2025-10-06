package handlers

import (
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/mappers"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OverviewHandler handles HTTP requests for the overview endpoint.
type OverviewHandler struct {
	service services.OverviewService
	logger  *zap.Logger
}

// NewOverviewHandler creates a new OverviewHandler.
func NewOverviewHandler(service services.OverviewService, logger *zap.Logger) *OverviewHandler {
	return &OverviewHandler{
		service: service,
		logger:  logger,
	}
}

// Get handles the request for obtaining the inventory overview.
//
//	@Summary		Get inventory overview
//	@Description	Returns a comprehensive overview of the system's inventory.
//	@Tags			Overview
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.OverviewSummary
//	@Failure		500	{object}	dto.GenericErrorResponse
//	@Router			/overview [get]
func (h *OverviewHandler) Get(c *gin.Context) {
	overview, err := h.service.GetOverview(c.Request.Context())
	if err != nil {
		h.logger.Error("error getting overview info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve overview: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, mappers.ToOverviewSummaryDTO(overview))
}
