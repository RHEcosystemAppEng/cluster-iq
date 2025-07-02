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
	APIStatus string `json:"api_status"`
	DBStatus  string `json:"db_status"`
}

// Check handles the request for checking the health level of the API.
//
//	@Summary		Runs HealthChecks
//	@Description	Runs several checks for evaluating the health level of ClusterIQ
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	healthCheckResponse
//	@Router			/health [get]
func (h *HealthCheckHandler) Check(c *gin.Context) {
	dbStatus := "OK"
	if err := h.db.Ping(); err != nil {
		h.logger.Error("Database health check failed", zap.Error(err))
		dbStatus = "Unavailable"
	}

	c.JSON(http.StatusOK, healthCheckResponse{
		APIStatus: "OK",
		DBStatus:  dbStatus,
	})
}
