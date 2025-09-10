package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// HealthCheckHandler handles health check requests.
type HealthCheckHandler struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewHealthCheckHandler creates a new HealthCheckHandler.
func NewHealthCheckHandler(db *sqlx.DB, logger *zap.Logger) *HealthCheckHandler {
	return &HealthCheckHandler{
		db:     db,
		logger: logger,
	}
}

type healthCheckResponse struct {
	APIHealth bool `json:"api_health"`
	DBHealth  bool `json:"db_health"`
}

// Check handles the request for checking the health of the API.
//
//	@Summary		Runs HealthChecks
//	@Description	Runs checks to evaluate the health of ClusterIQ.
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	healthCheckResponse
//	@Router			/healthcheck [get]
func (h *HealthCheckHandler) Check(c *gin.Context) {
	dbStatus := false
	if err := h.db.Ping(); err == nil {
		dbStatus = true
	} else {
		h.logger.Error("Database health check failed", zap.Error(err))
	}

	c.JSON(http.StatusOK, healthCheckResponse{
		APIHealth: true,
		DBHealth:  dbStatus,
	})
}
