package handlers

import (
	"net/http"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthCheckHandler exposes liveness/readiness-style endpoints.
type HealthCheckHandler struct {
	db     *dbclient.DBClient
	logger *zap.Logger
}

// NewHealthCheckHandler wires DB client and logger into the handler.
func NewHealthCheckHandler(db *dbclient.DBClient, logger *zap.Logger) *HealthCheckHandler {
	return &HealthCheckHandler{
		db:     db,
		logger: logger,
	}
}

// healthCheckResponse defines the response output format for the healthcheck endpoint
type healthCheckResponse struct {
	APIHealth bool `json:"api_health"`
	DBHealth  bool `json:"db_health"`
} // @name HealthCheckResponse

// Check returns current API and database health.
//
//	@Summary		Health checks
//	@Description	Report API process health and DB connectivity.
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	healthCheckResponse
//	@Router			/healthcheck [get]
func (h *HealthCheckHandler) Check(c *gin.Context) {
	dbStatus := false

	if err := h.db.Ping(); err != nil {
		h.logger.Error("Database health check failed", zap.Error(err))
	} else {
		dbStatus = true
	}

	c.JSON(http.StatusOK, healthCheckResponse{
		APIHealth: true,
		DBHealth:  dbStatus,
	})
}
