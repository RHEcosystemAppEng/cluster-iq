package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ==================== Health Checks Handlers ====================

// HandlerHealthCheck handles the request for checking the health level of the API
//
//	@Summary		Runs HealthChecks
//	@Description	Runs several checks for evaluating the health level of ClusterIQ
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthCheckResponse
//	@Router			/healthcheck [get]
func (a APIServer) HandlerHealthCheck(c *gin.Context) {
	hc := HealthChecks{
		APIHealth: false,
		DBHealth:  false,
	}

	// Checking DB Connection status
	if err := a.sql.Ping(); err == nil {
		hc.DBHealth = true
	} else {
		a.logger.Error("Can't ping DB", zap.Error(err))
	}

	// Checking API's Router status
	if a.router != nil {
		hc.APIHealth = true
	}

	c.PureJSON(http.StatusOK, HealthCheckResponse{HealthChecks: hc})
}
