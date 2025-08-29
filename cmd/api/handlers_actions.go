package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/events"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ==================== Scheduled Actions Handlers ====================

// HandlerGetScheduledActions retrieves all scheduled actions with optional filtering
//
//	@Summary		List all scheduled actions
//	@Description	Returns a list of scheduled actions
//	@Tags			Actions
//	@Param			enabled	query		string	false	"Filter by enabled state (true/false)"
//	@Param			status	query		string	false	"Filter by action status"
//	@Success		200		{object}	ScheduledActionListResponse
//	@Failure		500		{object}	GenericErrorResponse
//	@Router			/schedule [get]
func (a APIServer) HandlerGetScheduledActions(c *gin.Context) {
	a.logger.Debug("Retrieving scheduled actions")

	// Capturing query params
	var conditions []string
	var args []interface{}

	enabled := c.Query("enabled")
	if enabled != "" {
		conditions = append(conditions, "schedule.enabled = ?")
		args = append(args, enabled)
	}

	status := c.Query("status")
	if status != "" {
		conditions = append(conditions, "schedule.status = ?")
		args = append(args, status)
	}

	// Running sql client function
	schedule, err := a.sql.GetScheduledActions(conditions, args)
	if err != nil {
		a.logger.Error("Failed to retrieve scheduled actions", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewScheduledActionListResponse(schedule))
}

// HandlerGetScheduledActionByID retrieves a single scheduled action by its unique identifier
//
//	@Summary		Get scheduled action by ID
//	@Description	Returns details of a specific scheduled action identified by the action_id parameter
//	@Tags			Actions
//	@Param			action_id	path		string	true	"Scheduled action identifier"
//	@Success		200			{object}	ScheduledActionListResponse
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/schedule/{action_id} [get]
func (a APIServer) HandlerGetScheduledActionByID(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Retrieving scheduled action by ID", zap.String("action_id", actionID))

	schedule, err := a.sql.GetScheduledActionByID(actionID)
	if err != nil {
		a.logger.Error("Failed to retrieve scheduled action", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewScheduledActionListResponse(schedule))
}

// HandlerEnableScheduledAction activates a scheduled action so it can be executed according to its schedule
//
//	@Summary		Enable scheduled action
//	@Description	Activates a scheduled action specified by action_id
//	@Tags			Actions
//	@Param			action_id	path		string	true	"Scheduled action identifier"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/schedule/{action_id}/enable [patch]
func (a APIServer) HandlerEnableScheduledAction(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Enabling scheduled action", zap.String("action_id", actionID))

	err := a.sql.EnableScheduledAction(actionID)
	if err != nil {
		a.logger.Error("Failed to enable scheduled action", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerDisableScheduledAction deactivates a scheduled action to prevent its execution
//
//	@Summary		Disable scheduled action
//	@Description	Deactivates a scheduled action specified by action_id
//	@Tags			Actions
//	@Param			action_id	path		string	true	"Scheduled action identifier"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/schedule/{action_id}/disable [patch]
func (a APIServer) HandlerDisableScheduledAction(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Disabling action", zap.String("action_id", actionID))

	err := a.sql.DisableScheduledAction(actionID)
	if err != nil {
		a.logger.Error("Failed to disable scheduled action", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerPostScheduledAction processes the creation of new scheduled actions
//
//	@Summary		Create scheduled actions
//	@Description	Creates and registers new scheduled actions
//	@Tags			Actions
//	@Param			actions	body		[]json.RawMessage	true	"Scheduled actions to create"
//	@Success		200		{object}	nil
//	@Failure		500		{object}	GenericErrorResponse
//	@Router			/schedule [post]
func (a APIServer) HandlerPostScheduledAction(c *gin.Context) {
	a.logger.Debug("Creating scheduled actions")

	// Getting scheduled actions list on request's body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Var for Unmarshalling results
	var result []json.RawMessage

	// Unmarshalling response
	err = json.Unmarshal(body, &result)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Unmarshalling Actions by type
	decodedActions, err := actions.DecodeActions(result)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Writing scheduled action
	a.logger.Debug("Writing a new Scheduled Action", zap.Reflect("actions", decodedActions))
	err = a.sql.WriteScheduledActions(*decodedActions)
	if err != nil {
		a.logger.Error("Failed to create scheduled actions", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// TODO
	// We should return at least ID, nil is not useful
	c.PureJSON(http.StatusOK, nil)
}

// HandlerPatchStatusScheduledActions modifies only the status field of a scheduled action
//
//	@Summary		Update scheduled action status
//	@Description	Updates only the status field of a specific scheduled action identified by action_id
//	@Tags			Actions
//	@Param			action_id	path		string	true	"Scheduled action identifier"
//	@Param			status		query		string	true	"New status value"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	GenericErrorResponse	"Invalid input"
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/schedule/{action_id}/status [patch]
func (a APIServer) HandlerPatchStatusScheduledActions(c *gin.Context) {
	a.logger.Debug("Patching status of Scheduled Action Status")

	actionID := c.Param("action_id")
	status := c.Query("status")
	if status == "" {
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse("Status parameter is required"))
		return
	}

	// Writing scheduled action
	err := a.sql.PatchScheduledActionStatus(actionID, status)
	if err != nil {
		a.logger.Error("Failed to update scheduled action status", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerPatchScheduledActions processes updates to scheduled actions
//
//	@Summary		Update scheduled actions
//	@Description	Updates multiple fields of scheduled actions
//	@Tags			Actions
//	@Param			actions	body		[]json.RawMessage	true	"Scheduled actions to update"
//	@Success		200		{object}	nil
//	@Failure		500		{object}	GenericErrorResponse
//	@Router			/schedule [patch]
func (a APIServer) HandlerPatchScheduledActions(c *gin.Context) {
	a.logger.Debug("Updating scheduled actions")

	// Getting scheduled actions list on request's body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Var for Unmarshalling results
	var result []json.RawMessage

	// Unmarshalling response
	err = json.Unmarshal(body, &result)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Unmarshalling Actions by type
	decodedActions, err := actions.DecodeActions(result)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Writing scheduled action
	a.logger.Debug("Patching Scheduled Actions", zap.Int("action_count", len(*decodedActions)))
	err = a.sql.PatchScheduledAction(*decodedActions)
	if err != nil {
		a.logger.Error("Failed to update scheduled actions", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerDeleteScheduledAction permanently removes a scheduled action
//
//	@Summary		Delete scheduled action
//	@Description	Permanently removes a scheduled action identified by action_id
//	@Tags			Actions
//	@Param			action_id	path		string	true	"Scheduled action identifier"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/schedule/{action_id} [delete]
func (a APIServer) HandlerDeleteScheduledAction(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Removing a Scheduled Action", zap.String("action_id", actionID))

	if err := a.sql.DeleteScheduledAction(actionID); err != nil {
		a.logger.Error("Failed to delete scheduled action", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerPowerOnCluster handles startup of cluster instances
//
//	@Summary		Power on cluster
//	@Description	Starts all instances in the specified cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id}/power_on [post]
func (a APIServer) HandlerPowerOnCluster(c *gin.Context) {
	// TODO. We must add validation logic (middleware, validator, whatever)
	clusterID := c.Param("cluster_id")

	var request struct {
		TriggeredBy string  `json:"triggered_by"`
		Description *string `json:"description,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse("Invalid request body"))
		return
	}

	a.logger.Debug("Power On Cluster request received",
		zap.String("cluster_id", clusterID),
		zap.String("triggered_by", request.TriggeredBy))

	resp, err := a.handlePowerOn(clusterID, request.TriggeredBy, request.Description)
	if err != nil {
		a.logger.Error("Failed to power on cluster",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, resp)
}

func (a APIServer) handlePowerOn(clusterID, triggeredBy string, description *string) (*ClusterStatusChangeResponse, error) {
	a.logger.Debug("Powering On Cluster",
		zap.String("cluster_id", clusterID),
		zap.String("triggered_by", triggeredBy))

	// Initialize event tracker
	tracker := a.eventService.StartTracking(&events.EventOptions{
		Action:       inventory.ClusterPowerOnAction,
		Description:  description,
		ResourceID:   clusterID,
		ResourceType: inventory.ClusterResourceType,
		Result:       events.ResultPending,
		Severity:     events.SeverityInfo,
		TriggeredBy:  triggeredBy,
	})

	// Getting a new ClusterStatusChangeRequest for building the gRPC request
	cscr, err := NewClusterStatusChangeRequest(a.sql, clusterID)
	if err != nil {
		a.logger.Error("Cannot get ClusterStatusChangeRequest for the PowerOn gRPC request",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		tracker.Failed()
		return nil, fmt.Errorf("cannot get cluster status: %w", err)
	}

	// RPC call for power on a cluster
	if err := a.grpc.PowerOnCluster(cscr); err != nil {
		a.logger.Error("Error processing Cluster Power On request",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		tracker.Failed()
		return nil, fmt.Errorf("error processing power on request: %w", err)
	}

	a.logger.Info("Cluster Powered On successfully", zap.String("cluster_id", clusterID))

	// Update cluster status in DB
	if err := a.sql.UpdateClusterStatusByClusterID("Running", clusterID); err != nil {
		a.logger.Error("Error updating status in DB",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		tracker.Failed()
		return nil, fmt.Errorf("error updating cluster status: %w", err)
	}

	// Log successful completion
	tracker.Success()

	return NewClusterStatusChangeResponse(
		cscr.AccountName,
		cscr.ClusterID,
		cscr.Region,
		inventory.Running,
		cscr.InstancesIdList,
		nil,
	), nil
}

// HandlerPowerOffCluster handles graceful shutdown of cluster instances
//
//	@Summary		Power off cluster
//	@Description	Gracefully stops all instances in the specified cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id}/power_off [post]
func (a APIServer) HandlerPowerOffCluster(c *gin.Context) {
	// TODO. We must add validation logic (middleware, validator, whatever)
	clusterID := c.Param("cluster_id")

	var request struct {
		TriggeredBy string  `json:"triggered_by"`
		Description *string `json:"description,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse("Invalid request body"))
		return
	}

	a.logger.Debug("Power Off Cluster request received",
		zap.String("cluster_id", clusterID),
		zap.String("triggered_by", request.TriggeredBy))

	resp, err := a.handlePowerOff(clusterID, request.TriggeredBy, request.Description)
	if err != nil {
		a.logger.Error("Failed to power off cluster",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, resp)
}

func (a APIServer) handlePowerOff(clusterID, triggeredBy string, description *string) (*ClusterStatusChangeResponse, error) {
	a.logger.Debug("Powering Off Cluster",
		zap.String("cluster_id", clusterID),
		zap.String("triggered_by", triggeredBy))

	// Initialize event tracker
	tracker := a.eventService.StartTracking(&events.EventOptions{
		Action:       inventory.ClusterPowerOffAction,
		Description:  description,
		ResourceID:   clusterID,
		ResourceType: inventory.ClusterResourceType,
		Result:       events.ResultPending,
		Severity:     events.SeverityWarning,
		TriggeredBy:  triggeredBy,
	})
	// Getting a new ClusterStatusChangeRequest for building the gRPC request
	cscr, err := NewClusterStatusChangeRequest(a.sql, clusterID)
	if err != nil {
		a.logger.Error("Cannot get ClusterStatusChangeRequest for the PowerOff gRPC request",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		tracker.Failed()
		return nil, fmt.Errorf("cannot get cluster status: %w", err)
	}

	// RPC call for power off a cluster
	if err := a.grpc.PowerOffCluster(cscr); err != nil {
		a.logger.Error("Error processing Cluster Power Off request",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		tracker.Failed()
		return nil, fmt.Errorf("error processing power off request: %w", err)
	}

	a.logger.Info("Cluster Powered Off successfully", zap.String("cluster_id", clusterID))

	// Update cluster status in DB
	if err := a.sql.UpdateClusterStatusByClusterID("Stopped", clusterID); err != nil {
		a.logger.Error("Error updating status in DB",
			zap.String("cluster_id", clusterID),
			zap.Error(err))
		tracker.Failed()
		return nil, fmt.Errorf("error updating cluster status: %w", err)
	}

	// Log successful completion
	tracker.Success()

	return NewClusterStatusChangeResponse(
		cscr.AccountName,
		cscr.ClusterID,
		cscr.Region,
		inventory.Stopped,
		cscr.InstancesIdList,
		nil,
	), nil
}
