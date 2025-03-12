// TODO: Review documentaiton comments
// TODO: Generic response instead of NIL
// TODO: REview functions order
package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/actions"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
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

// ==================== Sched Actions Handlers ====================

// HandlerGetScheduledActions handles the request to obtain the entire Scheduled Actions list. Filters can be added using query params
//
//	@Summary		  Obtain every Scheduled action
//	@Description	Returns a list of Scheduled actions
//	@Tags			    actions
//	@Accept			  json
//	@Produce		  json
//	@Success		  200	{object}	ScheduledActionListResponse
//	@Failure		  500	{object}	GenericErrorResponse
//	@Router			  /schedule [get]
func (a APIServer) HandlerGetScheduledActions(c *gin.Context) {
	a.logger.Debug("Retrieving schedule action list")

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
		a.logger.Error("Can't retrieve Schedule list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
	}

	c.PureJSON(http.StatusOK, NewScheduledActionListResponse(schedule))
}

// HandlerGetScheduledActionsByID handles the request to obtain a specific scheduled action
//
//	@Summary		  Obtain a specific Scheduled action
//	@Description	Returns a Scheduled actions
//	@Tags			    actions
//	@Accept			  json
//	@Produce		  json
//	@Success		  200	{object}	ScheduledActionListResponse
//	@Failure		  500	{object}	GenericErrorResponse
//	@Router			  /schedule [get]
func (a APIServer) HandlerGetScheduledActionByID(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Retrieving complete scheduled actions list")

	schedule, err := a.sql.GetScheduledActionByID(actionID)
	if err != nil {
		a.logger.Error("Can't retrieve Schedule Action", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
	}

	c.PureJSON(http.StatusOK, NewScheduledActionListResponse(schedule))
}

// HandlerEnableScheduledAction handles the request to enable a specific scheduled action
//
//	@Summary		  Enables a specific Scheduled action
//	@Description	Based on the action_id, enables the action to be executed based on its definition
//	@Tags			    actions
//	@Accept			  json
//	@Produce		  json
//	@Success		  200	{object}	nil
//	@Failure		  500	{object}	GenericErrorResponse
//	@Router			  /schedule [get]
func (a APIServer) HandlerEnableScheduledAction(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Enabling action", zap.String("action_id", actionID))

	err := a.sql.EnableScheduledAction(actionID)
	if err != nil {
		a.logger.Error("Can't enable Schedule Action", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerDisableScheduledAction handles the request to enable a specific scheduled action
//
//	@Summary		  Disables a specific Scheduled action
//	@Description	Based on the action_id, disables the action to be executed based on its definition
//	@Tags			    actions
//	@Accept			  json
//	@Produce		  json
//	@Success		  200	{object}	nil
//	@Failure		  500	{object}	GenericErrorResponse
//	@Router			  /schedule [get]
func (a APIServer) HandlerDisableScheduledAction(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Enabling action", zap.String("action_id", actionID))

	err := a.sql.DisableScheduledAction(actionID)
	if err != nil {
		a.logger.Error("Can't disable Schedule Action", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerPostScheduledAction handles the request to create a new Scheduled Action
//
//	@Summary		  Creates a new Scheduled action
//	@Description	Creates and registers a new scheduled action on the DB
//	@Tags			    actions
//	@Accept			  json
//	@Produce		  json
//	@Success		  200	{object}  nil
//	@Failure		  500	{object}	GenericErrorResponse
//	@Router			  /schedule [get]
func (a APIServer) HandlerPostScheduledAction(c *gin.Context) {
	a.logger.Debug("Inserting list of Scheduled Actions")
	// Getting scheduled actions list on request's body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	var actions []actions.ScheduledAction
	err = json.Unmarshal(body, &actions)
	if err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
		return
	}

	// Writing scheduled action
	a.logger.Debug("Writing a new Scheduled Action", zap.Reflect("actions", actions))
	err = a.sql.WriteScheduledActions(actions)
	if err != nil {
		a.logger.Error("Can't write new Scheduled Actions into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerDeleteScheduledAction handles the request to remove a Scheduled Action
//
//	@Summary		  Deletes a Scheduled action
//	@Description	Deletes a specified scheduled action from the DB
//	@Tags			    actions
//	@Accept			  json
//	@Produce		  json
//	@Success		  200	{object}	nil
//	@Failure		  500	{object}	GenericErrorResponse
//	@Router			  /schedule [get]
func (a APIServer) HandlerDeleteScheduledAction(c *gin.Context) {
	actionID := c.Param("action_id")
	a.logger.Debug("Removing an Scheduled Action", zap.String("action_id", actionID))

	if err := a.sql.DeleteScheduledAction(actionID); err != nil {
		a.logger.Error("Can't delete instance from DB", zap.String("action_id", actionID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// ==================== Expenses      Handlers ====================

// HandlerGetExpenses handles the request for obtain the entire Expenses list
//
//	@Summary		Obtain every Expense
//	@Description	Returns a list of Expenses with every expense in the inventory
//	@Tags			Expenses
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ExpenseListResponse
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/expenses [get]
func (a APIServer) HandlerGetExpenses(c *gin.Context) {
	a.logger.Debug("Retrieving complete expenses list")

	expenses, err := a.sql.GetExpenses()
	if err != nil {
		a.logger.Error("Can't retrieve Expenses list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewExpenseListResponse(expenses))
}

// HandlerGetExpensesByInstance HandlerGetExpenseByID handles the request for obtain an Expense by its ID
//
//	@Summary		Obtain a single Expense by its ID
//	@Description	Returns a list of Expenses with a single Expense filtered by ID
//	@Tags			Expenses
//	@Accept			json
//	@Produce		json
//	@Param			instance_id	path		string	true	"Instance ID"
//	@Success		200			{object}	ExpenseListResponse
//	@Failure		404			{object}	nil
//	@Router			/expenses/{instance_id} [get]
func (a APIServer) HandlerGetExpensesByInstance(c *gin.Context) {
	instanceID := c.Param("instance_id")
	a.logger.Debug("Retrieving expenses by InstanceID", zap.String("instance_id", instanceID))

	expenses, err := a.sql.GetExpensesByInstance(instanceID)
	if err != nil {
		a.logger.Error("Instance not found", zap.String("instance_id", instanceID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	c.PureJSON(http.StatusOK, NewExpenseListResponse(expenses))
}

// HandlerPostExpense handles the request for writing a new Expense in the inventory
//
//	@Summary		Creates a new Expense in the inventory
//	@Description	Receives and write into the DB the information for a new Expense
//	@Tags			Expenses
//	@Accept			json
//	@Produce		json
//	@Param			instance	body		[]inventory.Expense	true	"New Expense to be added"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	GenericErrorResponse
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/expenses [post]
func (a APIServer) HandlerPostExpense(c *gin.Context) {
	// Getting expenses list on request's body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	var expenses []inventory.Expense
	err = json.Unmarshal(body, &expenses)
	if err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
		return
	}

	// Writing expenses
	a.logger.Debug("Writing a new Expense", zap.Reflect("expenses", expenses))
	err = a.sql.WriteExpenses(expenses)
	if err != nil {
		a.logger.Error("Can't write new Expenses into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// ==================== Instances     Handlers ====================

// HandlerGetInstances handles the request for obtain the entire Instances list
//
//	@Summary		Obtain every Instance
//	@Description	Returns a list of Instances with every Instance in the inventory
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	InstanceListResponse
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/instances [get]
func (a APIServer) HandlerGetInstances(c *gin.Context) {
	a.logger.Debug("Retrieving complete instance inventory")

	instances, err := a.sql.GetInstances()
	if err != nil {
		a.logger.Error("Can't retrieve Instances list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewInstanceListResponse(instances))
}

// HandlerGetInstancesForBillingUpdate handles the request for obtain a list of instances that needs to update its billing information
//
//	@Summary		Obtain instances list with missing billing data
//	@Description	Returns a list of Instances with outdated expenses or without any expense
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	InstanceListResponse
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/instances/expense_update [get]
func (a APIServer) HandlerGetInstancesForBillingUpdate(c *gin.Context) {
	a.logger.Debug("Retrieving instances with outdated billing information")

	instances, err := a.sql.GetInstancesOutdatedBilling()
	if err != nil {
		a.logger.Error("Can't retrieve Last Expenses list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewInstanceListResponse(instances))
}

// HandlerGetInstanceByID handles the request for obtain an Instance by its ID
//
//	@Summary		Obtain a single Instance by its ID
//	@Description	Returns a list of Instances with a single Instance filtered by ID
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			instance_id	path		string	true	"Instance ID"
//	@Success		200			{object}	InstanceListResponse
//	@Failure		404			{object}	nil
//	@Router			/instances/{instance_id} [get]
func (a APIServer) HandlerGetInstanceByID(c *gin.Context) {
	instanceID := c.Param("instance_id")
	a.logger.Debug("Retrieving instance by ID", zap.String("instance_id", instanceID))

	instances, err := a.sql.GetInstanceByID(instanceID)
	if err != nil {
		a.logger.Error("Instance not found", zap.String("instance_id", instanceID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	c.PureJSON(http.StatusOK, NewInstanceListResponse(instances))
}

// HandlerPostInstance handles the request for writing a new Instance in the inventory
//
//	@Summary		Creates a new Instance in the inventory
//	@Description	Receives and write into the DB the information for a new Instance
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			instance	body		[]inventory.Instance	true	"New Instance to be added"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	GenericErrorResponse
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/instances [post]
func (a APIServer) HandlerPostInstance(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	var instances []inventory.Instance
	err = json.Unmarshal(body, &instances)
	if err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
		return
	}

	a.logger.Debug("Writing a new Instance", zap.Reflect("instance", instances))
	err = a.sql.WriteInstances(instances)
	if err != nil {
		a.logger.Error("Can't write new instances into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}
	c.PureJSON(http.StatusOK, nil)
}

// HandlerDeleteInstance handles the request for removing an Instance in the inventory
//
//	@Summary		Deletes an Instance in the inventory
//	@Description	Deletes an Instance present in the inventory by its ID
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			instance_id	path		string	true	"Instance ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Router			/instances/{instance_id} [delete]
//
// TODO: Not Implemented
func (a APIServer) HandlerDeleteInstance(c *gin.Context) {
	instanceID := c.Param("instance_id")
	a.logger.Debug("Removing an Instance", zap.String("instance_id", instanceID))

	if err := a.sql.DeleteInstance(instanceID); err != nil {
		a.logger.Error("Can't delete instance from DB", zap.String("instance_id", instanceID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerPatchInstance handles the request for patching an Instance in the inventory
//
//	@Summary		Patches an Instance in the inventory
//	@Description	Receives and patch into the DB the information for an existing Instance
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Param			instance	body		inventory.Instance	true	"Instance to be modified"
//	@Param			instance_id	path		string				true	"Instance ID"
//	@Failure		501			{object}	nil					"Not Implemented"
//	@Router			/instances/{instance_id} [patch]
//
// TODO: NOT IMPLEMENTED
func (a APIServer) HandlerPatchInstance(c *gin.Context) {
	instanceID := c.Param("instance_id")
	a.logger.Debug("Patching an Instance", zap.String("instance_id", instanceID))

	c.PureJSON(http.StatusNotImplemented, nil)
}

// ==================== Clusters      Handlers ====================

// HandlerGetClusters handles the request for obtaining the entire Cluster list
//
//	@Summary		Obtain every Cluster
//	@Description	Returns a list of Clusters with a single instance filtered by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ClusterListResponse
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/clusters [get]
func (a APIServer) HandlerGetClusters(c *gin.Context) {
	a.logger.Debug("Retrieving complete clusters inventory")

	clusters, err := a.sql.GetClusters()
	if err != nil {
		a.logger.Error("Can't retrieve Clusters list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewClusterListResponse(clusters))
}

// HandlerGetClustersByID handles the request for obtain a Cluster by its Name
//
//	@Summary		Obtain a single Cluster by its Name
//	@Description	Returns a list of Clusters with a single Cluster filtered by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	ClusterListResponse
//	@Failure		404			{object}	nil
//	@Router			/clusters/{cluster_id} [get]
func (a APIServer) HandlerGetClustersByID(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Retrieving Cluster Tags by ID", zap.String("cluster_id", clusterID))

	clusters, err := a.sql.GetClusterByID(clusterID)
	if err != nil {
		a.logger.Error("Cluster not found", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	c.PureJSON(http.StatusOK, NewClusterListResponse(clusters))
}

// HandlerGetInstancesOnCluster handles the request for obtain the list of Instances belonging to a specific Cluster
//
//	@Summary		Obtain Instances list belonging to a Cluster
//	@Description	Returns a list of Instances belonging to a Cluster given by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	InstanceListResponse
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id}/instances [get]
func (a APIServer) HandlerGetInstancesOnCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Retrieving Cluster's Instances", zap.String("cluster_id", clusterID))

	instances, err := a.sql.GetInstancesOnCluster(clusterID)
	if err != nil {
		a.logger.Error("Can't retrieve instances on cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewInstanceListResponse(instances))
}

// HandlerGetClusterTags handles the request for obtain the list of tags of a Cluster
//
//	@Summary		Obtain Cluster Tags
//	@Description	Returns a list of Tags belonging to a Cluster given by ID
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	TagListResponse
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id}/tags [get]
func (a APIServer) HandlerGetClusterTags(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Retrieving Cluster's Tags", zap.String("cluster_id", clusterID))

	tags, err := a.sql.GetClusterTags(clusterID)
	if err != nil {
		a.logger.Error("Can't retrieve Tags of cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewTagListResponse(tags))
}

// HandlerPostCluster handles the request for writing a new Cluster in the inventory
//
//	@Summary		Creates a new Cluster in the inventory
//	@Description	Receives and write into the DB the information for a new Cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster	body		inventory.Cluster	true	"New Cluster to be added"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	GenericErrorResponse
//	@Failure		500		{object}	GenericErrorResponse
//	@Router			/clusters [post]
func (a APIServer) HandlerPostCluster(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	var clusters []inventory.Cluster
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
		return
	}

	a.logger.Debug("Writing new Clusters", zap.Reflect("clusters", clusters))
	err = a.sql.WriteClusters(clusters)
	if err != nil {
		a.logger.Error("Can't write new Clusters into DB", zap.Error(err))
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
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Powering On Cluster", zap.String("cluster_id", clusterID))

	// Getting a new ClusterStatusChangeRequest for building the gRPC request
	cscr, err := NewClusterStatusChangeRequest(a.sql, clusterID)
	if err != nil {
		a.logger.Error("Cannot get ClusterStatusChangeRequest for the PowerOn gRPC request", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// RPC call for power on a cluster
	if err := a.grpc.PowerOnCluster(cscr); err != nil {
		a.logger.Error("Error processing Cluster Power On request", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}
	a.logger.Info("Cluster Powered On successfully", zap.String("cluster_id", clusterID))

	// Updating Cluster Status on the DB
	if err := a.sql.UpdateClusterStatusByClusterID("Running", clusterID); err != nil {
		a.logger.Error("Error updating status on DB when powering on a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK,
		NewClusterStatusChangeResponse(
			cscr.AccountName,
			cscr.ClusterID,
			cscr.Region,
			inventory.Running,
			cscr.InstancesIdList,
			nil,
		),
	)
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
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Powering Off Cluster", zap.String("cluster_id", clusterID))

	// Getting a new ClusterStatusChangeRequest for building the gRPC request
	cscr, err := NewClusterStatusChangeRequest(a.sql, clusterID)
	if err != nil {
		a.logger.Error("Cannot get ClusterStatusChangeRequest for the PowerOff gRPC request", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// RPC call for power off a cluster
	if err := a.grpc.PowerOffCluster(cscr); err != nil {
		a.logger.Error("Error processing Cluster Power Off request", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}
	a.logger.Info("Cluster Powered Off successfully", zap.String("cluster_id", clusterID))

	// Updating Cluster Status on the DB
	if err := a.sql.UpdateClusterStatusByClusterID("Stopped", clusterID); err != nil {
		a.logger.Error("Error updating status on DB when powering off a cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK,
		NewClusterStatusChangeResponse(
			cscr.AccountName,
			cscr.ClusterID,
			cscr.Region,
			inventory.Running,
			cscr.InstancesIdList,
			nil,
		),
	)
}

// HandlerDeleteCluster handles the request for removing a Cluster in the inventory
//
//	@Summary		Deletes a Cluster in the inventory
//	@Description	Deletes a Cluster present in the inventory by its Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id} [delete]
func (a APIServer) HandlerDeleteCluster(c *gin.Context) {
	clusterName := c.Param("cluster_id")
	a.logger.Debug("Removing a Cluster", zap.String("cluster_id", clusterName))

	if err := a.sql.DeleteCluster(clusterName); err != nil {
		a.logger.Error("Can't delete Cluster from DB", zap.String("cluster_id", clusterName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerPatchCluster handles the request for patching a Cluster in the inventory
//
//	@Summary		Patches a Cluster in the inventory
//	@Description	Receives and patch into the DB the information for an existing Cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string				true	"Cluster ID"
//	@Param			cluster		body		inventory.Cluster	true	"Cluster to be modified"
//	@Success		200			{object}	nil
//	@Failure		501			{object}	nil	"Not Implemented"
//	@Router			/clusters/{cluster_id} [patch]
//
// TODO: NOT IMPLEMENTED
func (a APIServer) HandlerPatchCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	a.logger.Debug("Patching a Cluster", zap.String("cluster_id", clusterID))

	c.PureJSON(http.StatusNotImplemented, nil)
}

// ==================== Accounts      Handlers ====================

// HandlerGetAccounts handles the request for obtaining the entire Account list
//
//	@Summary		Obtain every Account
//	@Description	Returns a list of Accounts with a single Account filtered by Name
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	AccountListResponse
//	@Failure		500	{object}	nil
//	@Router			/accounts [get]
func (a APIServer) HandlerGetAccounts(c *gin.Context) {
	a.logger.Debug("Retrieving complete Accounts inventory")

	accounts, err := a.sql.GetAccounts()
	if err != nil {
		a.logger.Error("Can't retrieve Accounts list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewAccountListResponse(accounts))
}

// HandlerGetAccountsByName handles the request for obtain an Account by its Name
//
//	@Summary		Obtain a single Account by its Name
//	@Description	Returns a list of Accounts with a single Account filtered by Name
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account_name	path		string	true	"Account Name"
//	@Success		200				{object}	AccountListResponse
//	@Failure		404				{object}	GenericErrorResponse
//	@Router			/accounts/{account_name} [get]
func (a APIServer) HandlerGetAccountsByName(c *gin.Context) {
	accountName := c.Param("account_name")
	a.logger.Debug("Retrieving Account by Name", zap.String("account_name", accountName))

	accounts, err := a.sql.GetAccountByName(accountName)
	if err != nil {
		a.logger.Error("Account not found", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusNotFound, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewAccountListResponse(accounts))
}

// HandlerGetClustersOnAccount handles the request for obtain the list of clusters deployed on a specific Account
//
//	@Summary		Obtain Cluster list on an Account
//	@Description	Returns a list of Clusters which belongs to an Account given by Name
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account_name	path		string	true	"Account Name"
//	@Success		200				{object}	ClusterListResponse
//	@Failure		500				{object}	nil
//	@Router			/accounts/{account_name}/clusters [get]
func (a APIServer) HandlerGetClustersOnAccount(c *gin.Context) {
	accountName := c.Param("account_name")
	a.logger.Debug("Retrieving Account's Clusters", zap.String("account_name", accountName))

	clusters, err := a.sql.GetClustersOnAccount(accountName)
	if err != nil {
		a.logger.Error("Can't retrieve clusters on account", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, NewClusterListResponse(clusters))
}

// HandlerPostAccount handles the request for writing a new Account in the inventory
//
//	@Summary		Creates a new Account in the inventory
//	@Description	Receives and write into the DB the information for a new Account
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account	body		inventory.Account	true	"New Account to be added"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	nil
//	@Failure		500		{object}	GenericErrorResponse
//	@Router			/accounts [post]
func (a APIServer) HandlerPostAccount(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	var accounts []inventory.Account
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, nil)
		return
	}

	a.logger.Debug("Writing a new Account", zap.Reflect("accounts", accounts))
	err = a.sql.WriteAccounts(accounts)
	if err != nil {
		a.logger.Error("Can't write new Accounts into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerDeleteAccount handles the request for deleting an Account in the inventory
//
//	@Summary		Deletes an Account in the inventory
//	@Description	Deletes an Account present in the inventory by its Name
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account_name	path		string	true	"Account Name"
//	@Success		200				{object}	nil
//	@Failure		500				{object}	GenericErrorResponse
//	@Router			/accounts/{account_name} [delete]
func (a APIServer) HandlerDeleteAccount(c *gin.Context) {
	accountName := c.Param("account_name")
	a.logger.Debug("Removing an Account", zap.String("account", accountName))

	if err := a.sql.DeleteAccount(accountName); err != nil {
		a.logger.Error("Can't delete Cluster from DB", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, nil)
}

// HandlerPatchAccount handles the request for patching an Account in the inventory
//
//	@Summary		Patches an Account in the inventory
//	@Description	Receives and patch into the DB the information for an existing Account
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			Account			body		inventory.Account	true	"Account to be modified"
//	@Param			account_name	path		string				true	"Account Name"
//	@Failure		501				{object}	nil					"Not Implemented"
//	@Router			/accounts/{account_name} [patch]
func (a APIServer) HandlerPatchAccount(c *gin.Context) {
	accountName := c.Param("account_name")
	a.logger.Debug("Patching an Account", zap.String("account", accountName))

	c.PureJSON(http.StatusNotImplemented, nil)
}

// ==================== Extra      Handlers ====================

// HandlerRefreshInventory handles the request for refreshing the entire
// inventory just after a full scan. This method is used for recalculating some
// values and mark the missing clusters as "terminated"
//
//	@Summary		Refresh data on inventory
//	@Description	Recalculating some values and mark the missing clusters as "terminated"
//	@Tags			Inventory
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	nil
//	@Failure		500	{object}	nil
//	@Router			/inventory/refresh [post]
func (a APIServer) HandlerRefreshInventory(c *gin.Context) {
	if err := a.sql.RefreshInventory(); err != nil {
		a.logger.Error("Can't refresh inventory data on DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}
	// This function doesn't return any 200OK code for preventing duplicated responses
}
