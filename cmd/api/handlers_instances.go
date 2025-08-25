package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	art "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	dtomodel "github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

	// Transforming into DTO type
	// TODO move to function
	var response []dtomodel.InstanceDTOResponse
	for _, instance := range instances {
		response = append(response, *instance.ToInstanceDTOResponse())
	}

	c.PureJSON(http.StatusOK, dtomodel.NewInstanceDTOResponseList(response))
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
		// TODO Fix this error message. Not always is "NotFound"
		// TODO Check also this in the rest of handlers
		a.logger.Error("Instance not found", zap.String("instance_id", instanceID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, nil)
		return
	}

	c.PureJSON(http.StatusOK, dtomodel.NewInstanceDTOResponseList([]dtomodel.InstanceDTOResponse{*instances.ToInstanceDTOResponse()}))
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

	c.PureJSON(http.StatusOK, instances)
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

	var instances dtomodel.InstanceDTORequestList
	if err = json.Unmarshal(body, &instances); err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
		return
	}

	a.logger.Debug("Writing a new Instance", zap.Reflect("instance", instances))

	// Filling Account internal ID for every instance
	var toWriteInstances []inventory.Instance
	for _, instance := range instances.Instances {
		newInstance := *instance.ToInventoryInstance()
		if id, err := a.sql.GetClusterInternalID(instance.ClusterID); err != nil {
			a.logger.Error("Can't obtain internal ID for Cluster", zap.Error(err))
			c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
			return
		} else {
			newInstance.ClusterID = id
			toWriteInstances = append(toWriteInstances, newInstance)
		}
	}

	if err = a.sql.WriteInstances(toWriteInstances); err != nil {
		a.logger.Error("Can't write new instances into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, art.PostResponse{Count: len(toWriteInstances), Status: "Instance(s) Post OK"})
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

	c.PureJSON(http.StatusOK, art.DeleteResponse{
		Count:  1,
		Status: fmt.Sprintf("Instance '%s' Delete OK", instanceID),
	})
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
