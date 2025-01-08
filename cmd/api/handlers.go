package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerGetExpenses handles the request for obtain the entire Expenses list
//
//	@Summary		Obtain every Expense
//	@Description	Returns a list of Expenses with every expense in the inventory
//	@Tags			Expenses
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ExpenseListResponse
//	@Failure		500	{object}	nil
//	@Router			/expenses [get]
func HandlerGetExpenses(c *gin.Context) {
	logger.Debug("Retrieving complete expense inventory")
	addHeaders(c)

	expenses, err := getExpenses()
	if err != nil {
		logger.Error("Can't retrieve Expenses list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewExpenseListResponse(expenses)
	c.PureJSON(http.StatusOK, response)
}

// HandlerGetInstancesForBillingUpdate handles the request for obtain a list of instances that needs to update its billing information
//
//	@Summary		Obtain instances list with missing billing data
//	@Description	Returns a list of Instances with outdated expenses or without any expense
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	InstanceListResponse
//	@Failure		500	{object}	nil
//	@Router			/instances/expense_update [get]
func HandlerGetInstancesForBillingUpdate(c *gin.Context) {
	logger.Debug("Retrieving instances with outdated billing information")
	addHeaders(c)

	instances, err := getInstancesOutdatedBilling()
	if err != nil {
		logger.Error("Can't retrieve Last Expenses list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewInstanceListResponse(instances)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		500			{object}	nil
//	@Router			/expenses/{instance_id} [get]
func HandlerGetExpensesByInstance(c *gin.Context) {
	instanceID := c.Param("instance_id")
	logger.Debug("Retrieving expenses by InstanceID", zap.String("instance_id", instanceID))
	addHeaders(c)

	expenses, err := getExpensesByInstance(instanceID)
	if err != nil {
		logger.Error("Instance not found", zap.String("instance_id", instanceID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	response := NewExpenseListResponse(expenses)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		400			{object}	nil
//	@Failure		500			{object}	nil
//	@Router			/expenses [post]
func HandlerPostExpense(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	var expenses []inventory.Expense
	err = json.Unmarshal(body, &expenses)
	if err != nil {
		logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, nil)
		return
	}

	logger.Debug("Writing a new Expense", zap.Reflect("expenses", expenses))
	err = writeExpenses(expenses)
	if err != nil {
		logger.Error("Can't write new Expenses into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}
	c.PureJSON(http.StatusOK, nil)
}

// HandlerHealthCheck handles the request for checking the health level of the API
//
//	@Summary		Runs HealthChecks
//	@Description	Runs several checks for evaluating the health level of ClusterIQ
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	HealthCheckResponse
//	@Router			/healthcheck [get]
func HandlerHealthCheck(c *gin.Context) {
	logger.Debug("Running Health Checks")
	addHeaders(c)

	hc := HealthChecks{
		APIHealth: false,
		DBHealth:  false,
	}

	// Checking DB Connection status
	if db != nil {
		if err := db.Ping(); err == nil {
			hc.DBHealth = true
		} else {
			logger.Error("Can't ping DB", zap.Error(err))
		}
	}

	// Checking API's Router status
	if router != nil {
		hc.APIHealth = true
	}

	response := HealthCheckResponse{HealthChecks: hc}
	c.PureJSON(http.StatusOK, response)
}

// HandlerGetInstances handles the request for obtain the entire Instances list
//
//	@Summary		Obtain every Instance
//	@Description	Returns a list of Instances with every Instance in the inventory
//	@Tags			Instances
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	InstanceListResponse
//	@Failure		500	{object}	nil
//	@Router			/instances [get]
func HandlerGetInstances(c *gin.Context) {
	logger.Debug("Retrieving complete instance inventory")
	addHeaders(c)

	instances, err := getInstances()
	if err != nil {
		logger.Error("Can't retrieve Instances list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewInstanceListResponse(instances)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		500			{object}	nil
//	@Router			/instances/{instance_id} [get]
func HandlerGetInstanceByID(c *gin.Context) {
	instanceID := c.Param("instance_id")
	logger.Debug("Retrieving instance by ID", zap.String("instance_id", instanceID))
	addHeaders(c)

	instances, err := getInstanceByID(instanceID)
	if err != nil {
		logger.Error("Instance not found", zap.String("instance_id", instanceID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	response := NewInstanceListResponse(instances)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		500			{object}	nil
//	@Router			/instances [post]
func HandlerPostInstance(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	var instances []inventory.Instance
	err = json.Unmarshal(body, &instances)
	if err != nil {
		logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, nil)
		return
	}

	logger.Debug("Writing a new Instance", zap.Reflect("instance", instances))
	err = writeInstances(instances)
	if err != nil {
		logger.Error("Can't write new instances into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
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
func HandlerDeleteInstance(c *gin.Context) {
	instanceID := c.Param("instance_id")
	logger.Debug("Removing an Instance", zap.String("instance_id", instanceID))
	if err := deleteInstance(instanceID); err != nil {
		logger.Error("Can't delete instance from DB", zap.String("instance_id", instanceID), zap.Error(err))
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
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Failure		501			{object}	nil	"Not Implemented"
//	@Router			/instances/{instance_id} [patch]
func HandlerPatchInstance(c *gin.Context) {
	instanceID := c.Param("instance_id")
	logger.Debug("Patching an Instance", zap.String("instance_id", instanceID))
	c.PureJSON(http.StatusNotImplemented, nil)
}

// HandlerGetClusters handles the request for obtaining the entire Cluster list
//
//	@Summary		Obtain every Cluster
//	@Description	Returns a list of Clusters with a single instance filtered by Name
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ClusterListResponse
//	@Failure		500	{object}	nil
//	@Router			/clusters [get]
func HandlerGetClusters(c *gin.Context) {
	logger.Debug("Retrieving complete clusters inventory")
	addHeaders(c)

	clusters, err := getClusters()
	if err != nil {
		logger.Error("Can't retrieve Clusters list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewClusterListResponse(clusters)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		500			{object}	nil
//	@Router			/clusters/{cluster_id} [get]
func HandlerGetClustersByID(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	logger.Debug("Retrieving Cluster Tags by ID", zap.String("cluster_id", clusterID))
	addHeaders(c)

	clusters, err := getClusterByID(clusterID)
	if err != nil {
		logger.Error("Cluster not found", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	response := NewClusterListResponse(clusters)
	c.PureJSON(http.StatusOK, response)
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
func HandlerGetClusterTags(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	logger.Debug("Retrieving Cluster's Tags", zap.String("cluster_id", clusterID))
	addHeaders(c)

	tags, err := getClusterTags(clusterID)
	if err != nil {
		logger.Error("Can't retrieve Tags of cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewTagListResponse(tags)
	c.PureJSON(http.StatusOK, response)
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
func HandlerGetInstancesOnCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	logger.Debug("Retrieving Cluster's Instances", zap.String("cluster_id", clusterID))
	addHeaders(c)

	instances, err := getInstancesOnCluster(clusterID)
	if err != nil {
		logger.Error("Can't retrieve instances on cluster", zap.String("cluster_id", clusterID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewInstanceListResponse(instances)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		500		{object}	nil
//	@Router			/clusters [post]
func HandlerPostCluster(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	var clusters []inventory.Cluster
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, nil)
		return
	}

	logger.Debug("Writing new Clusters", zap.Reflect("clusters", clusters))
	err = writeClusters(clusters)
	if err != nil {
		logger.Error("Can't write new Clusters into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}
	c.PureJSON(http.StatusOK, nil)
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
//
// TODO: Not Implemented
func HandlerDeleteCluster(c *gin.Context) {
	clusterName := c.Param("cluster_id")
	logger.Debug("Removing a Cluster", zap.String("cluster_id", clusterName))
	if err := deleteCluster(clusterName); err != nil {
		logger.Error("Can't delete Cluster from DB", zap.String("cluster_id", clusterName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
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
//	@Failure		500			{object}	nil
//	@Failure		501			{object}	nil	"Not Implemented"
//	@Router			/clusters/{cluster_id} [patch]
func HandlerPatchCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	logger.Debug("Patching a Cluster", zap.String("cluster_id", clusterID))
	c.PureJSON(http.StatusNotImplemented, nil)
}

// HandlerStopCluster handles graceful shutdown of cluster instances
//
//	@Summary		Power off cluster
//	@Description	Gracefully stops all instances in the specified cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Failure		501			{object}	nil	"Not Implemented"
//	@Router			/clusters/{cluster_id}/power_off [post]
func HandlerStopCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	logger.Debug("Stopping the cluster", zap.String("cluster_id", clusterID))
	c.PureJSON(http.StatusNotImplemented, nil)
}

// HandlerStartCluster handles startup of cluster instances
//
//	@Summary		Power on cluster
//	@Description	Starts all instances in the specified cluster
//	@Tags			Clusters
//	@Accept			json
//	@Produce		json
//	@Param			cluster_id	path		string	true	"Cluster ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	nil
//	@Failure		501			{object}	nil	"Not Implemented"
//	@Router			/clusters/{cluster_id}/power_on [post]
func HandlerStartCluster(c *gin.Context) {
	clusterID := c.Param("cluster_id")
	logger.Debug("Starting the Cluster", zap.String("cluster_id", clusterID))
	c.PureJSON(http.StatusNotImplemented, nil)
}

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
func HandlerGetAccounts(c *gin.Context) {
	logger.Debug("Retrieving complete Accounts inventory")
	addHeaders(c)

	accounts, err := getAccounts()
	if err != nil {
		logger.Error("Can't retrieve Accounts list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewAccountListResponse(accounts)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		404				{object}	nil
//	@Failure		500				{object}	nil
//	@Router			/accounts/{account_name} [get]
func HandlerGetAccountsByName(c *gin.Context) {
	accountName := c.Param("account_name")
	logger.Debug("Retrieving Account by Name", zap.String("account_name", accountName))
	addHeaders(c)

	accounts, err := getAccountByName(accountName)
	if err != nil {
		logger.Error("Account not found", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	response := NewAccountListResponse(accounts)
	c.PureJSON(http.StatusOK, response)
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
func HandlerGetClustersOnAccount(c *gin.Context) {
	accountName := c.Param("account_name")
	logger.Debug("Retrieving Account's Clusters", zap.String("account_name", accountName))
	addHeaders(c)

	clusters, err := getClustersOnAccount(accountName)
	if err != nil {
		logger.Error("Can't retrieve clusters on account", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewClusterListResponse(clusters)
	c.PureJSON(http.StatusOK, response)
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
//	@Failure		500		{object}	nil
//	@Router			/accounts [post]
func HandlerPostAccount(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	var accounts []inventory.Account
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, nil)
		return
	}

	logger.Debug("Writing a new Account", zap.Reflect("accounts", accounts))
	err = writeAccounts(accounts)
	if err != nil {
		logger.Error("Can't write new Accounts into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
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
//	@Failure		500				{object}	nil
//	@Router			/accounts/{account_name} [delete]
//
// TODO: Not Implemented
func HandlerDeleteAccount(c *gin.Context) {
	accountName := c.Param("account_name")
	logger.Debug("Removing an Account", zap.String("account", accountName))
	if err := deleteAccount(accountName); err != nil {
		logger.Error("Can't delete Cluster from DB", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
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
//	@Success		200				{object}	nil
//	@Failure		500				{object}	nil
//	@Failure		501				{object}	nil	"Not Implemented"
//	@Router			/accounts/{account_name} [patch]
func HandlerPatchAccount(c *gin.Context) {
	accountName := c.Param("account_name")
	logger.Debug("Patching an Account", zap.String("account", accountName))
	c.PureJSON(http.StatusNotImplemented, nil)
}

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
func HandlerRefreshInventory(c *gin.Context) {
	if err := refreshInventory(); err != nil {
		logger.Error("Can't refresh inventory data on DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}
}
