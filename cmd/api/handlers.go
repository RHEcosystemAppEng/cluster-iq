package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerGetInstances handles /instances endpoint
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

// HandlerGetInstancesByID handles /instances/:instance_id endpoint
func HandlerGetInstancesByID(c *gin.Context) {
	logger.Debug("Retrieving instance by ID")
	addHeaders(c)

	instanceID := c.Param("instance_id")

	instances, err := getInstanceByID(instanceID)
	if err != nil {
		logger.Error("Instance not found", zap.String("instance_id", instanceID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	response := NewInstanceListResponse(instances)
	c.PureJSON(http.StatusOK, response)
}

// HandlerGetClusters handles /clusters endpoint
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

// HandlerGetClustersByName handles /clusters/:cluster_name endpoint
func HandlerGetClustersByName(c *gin.Context) {
	logger.Debug("Retrieving cluster by Name")
	addHeaders(c)

	clusterName := c.Param("cluster_name")

	clusters, err := getClusterByName(clusterName)
	if err != nil {
		logger.Error("Cluster not found", zap.String("cluster_name", clusterName), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	response := NewClusterListResponse(clusters)
	c.PureJSON(http.StatusOK, response)
}

// HandlerGetInstancesOnCluster handles /accounts/:account_name/cluster endpoint
func HandlerGetInstancesOnCluster(c *gin.Context) {
	logger.Debug("Retrieving Account's Clusters")
	addHeaders(c)

	clusterName := c.Param("cluster_name")

	instances, err := getInstancesOnCluster(clusterName)
	if err != nil {
		logger.Error("Can't retrieve instances on cluster", zap.String("cluster_name", clusterName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewInstanceListResponse(instances)
	c.PureJSON(http.StatusOK, response)
}

// HandlerGetAccounts handles /accounts/ endpoint
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

// HandlerGetAccountsByName handles /accounts/:account_name endpoint
func HandlerGetAccountsByName(c *gin.Context) {
	logger.Debug("Retrieving Account by Name")
	addHeaders(c)

	accountName := c.Param("account_name")

	accounts, err := getAccountByName(accountName)
	if err != nil {
		logger.Error("Account not found", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	response := NewAccountListResponse(accounts)
	c.PureJSON(http.StatusOK, response)
}

// HandlerGetClustersOnAccount handles /accounts/:account_name/cluster endpoint
func HandlerGetClustersOnAccount(c *gin.Context) {
	logger.Debug("Retrieving Account's Clusters")
	addHeaders(c)

	accountName := c.Param("account_name")

	clusters, err := getClustersOnAccount(accountName)
	if err != nil {
		logger.Error("Can't retrieve clusters on account", zap.String("account_name", accountName), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, nil)
		return
	}

	response := NewClusterListResponse(clusters)
	c.PureJSON(http.StatusOK, response)
}
