package handlers

import (
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InventoryHandler connects HTTP endpoints with the InventoryService.
type InventoryHandler struct {
	service services.InventoryService
	logger  *zap.Logger
}

// NewInventoryHandler builds a new InventoryHandler with its dependencies.
func NewInventoryHandler(service services.InventoryService, logger *zap.Logger) *InventoryHandler {
	return &InventoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *InventoryHandler) Refresh(c *gin.Context) {
	if err := h.service.Refresh(c); err != nil {
		h.logger.Error("Error when refreshing Inventory", zap.Error(err))
	}
}
