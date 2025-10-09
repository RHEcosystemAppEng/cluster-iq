package handlers

import (
	"net/http"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OverviewHandler exposes the inventory overview endpoint.
type OverviewHandler struct {
	service services.OverviewService
	logger  *zap.Logger
}

// NewOverviewHandler returns an OverviewHandler with its dependencies.
func NewOverviewHandler(service services.OverviewService, logger *zap.Logger) *OverviewHandler {
	return &OverviewHandler{
		service: service,
		logger:  logger,
	}
}

// Get returns the aggregated inventory overview.
//
//	@Summary		Get inventory overview
//	@Description	Return an aggregated summary of the inventory state.
//	@Tags			Overview
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.OverviewSummary
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/overview [get]
func (h *OverviewHandler) Get(c *gin.Context) {
	overview, err := h.service.GetOverview(c.Request.Context())
	if err != nil {
		h.logger.Error("error getting overview info", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve overview: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToOverviewSummaryDTO(overview))
}
