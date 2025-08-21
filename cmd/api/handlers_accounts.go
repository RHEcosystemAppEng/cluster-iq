package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/apiresponsetypes"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ==================== Accounts      Handlers ====================

// HandlerGetAccounts handles the request for obtaining the entire Account list
//
//	@Summary		Obtain every Account
//	@Description	Returns a list of Accounts with a single Account filtered by Name
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.AccountDTOResponseList
//	@Failure		500	{object}	GenericErrorResponse
//	@Router			/accounts [get]
func (a APIServer) HandlerGetAccounts(c *gin.Context) {
	a.logger.Debug("Retrieving complete Accounts inventory")

	accounts, err := a.sql.GetAccounts()
	if err != nil {
		a.logger.Error("Can't retrieve Accounts list", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Transforming into DTO type
	// TODO move to function
	var response []dto.AccountDTOResponse
	for _, account := range accounts {
		response = append(response, *account.ToAccountDTOResponse())
	}

	c.PureJSON(http.StatusOK, dto.NewAccountDTOResponseList(response))
}

// HandlerGetAccountsByID handles the request for obtain an Account by its Name
//
//	@Summary		Obtain a single Account by its Name
//	@Description	Returns a list of Accounts with a single Account filtered by Name
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account_id	path		string	true	"Account Name"
//	@Success		200			{object}	dto.AccountDTOResponseList
//	@Failure		404			{object}	GenericErrorResponse
//	@Router			/accounts/{account_id} [get]
func (a APIServer) HandlerGetAccountsByID(c *gin.Context) {
	accountID := c.Param("account_id")
	a.logger.Debug("Retrieving Account by ID", zap.String("account_id", accountID))

	accounts, err := a.sql.GetAccountByID(accountID)
	if err != nil {
		a.logger.Error("Account not found", zap.String("account_id", accountID), zap.Error(err))
		c.PureJSON(http.StatusNotFound, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, dto.NewAccountDTOResponseList([]dto.AccountDTOResponse{*accounts.ToAccountDTOResponse()}))
}

// HandlerGetClustersOnAccount handles the request for obtain the list of clusters deployed on a specific Account
//
//	@Summary		Obtain Cluster list on an Account
//	@Description	Returns a list of Clusters which belongs to an Account given by Name
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account_id	path		string	true	"Account Name"
//	@Success		200			{object}	ClusterListResponse
//	@Failure		500			{object}	nil
//	@Router			/accounts/{account_id}/clusters [get]
func (a APIServer) HandlerGetClustersOnAccount(c *gin.Context) {
	account_id := c.Param("account_id")
	a.logger.Debug("Retrieving Account's Clusters", zap.String("account_id", account_id))

	clusters, err := a.sql.GetClustersOnAccount(account_id)
	if err != nil {
		a.logger.Error("Can't retrieve clusters on account", zap.String("account_id", account_id), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	// Transforming into DTO type
	var response []dto.ClusterDTOResponse
	for _, cluster := range clusters {
		response = append(response, *cluster.ToClusterDTOResponse())
	}
	c.PureJSON(http.StatusOK, dto.NewClusterDTOResponseList(response))
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
	a.logger.Debug("Writing new Accounts")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		a.logger.Error("Can't get body from request", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	var accounts dto.AccountDTORequestList
	err = json.Unmarshal(body, &accounts)
	if err != nil {
		a.logger.Error("Can't obtain data from body request", zap.Error(err))
		c.PureJSON(http.StatusBadRequest, NewGenericErrorResponse(err.Error()))
		return
	}

	if err = a.sql.WriteAccounts(*accounts.ToInventoryAccountList()); err != nil {
		a.logger.Error("Can't write new Accounts into DB", zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, apiresponsetypes.PostResponse{Count: len(accounts.Accounts), Status: "Account(s) Post OK"})
}

// HandlerDeleteAccount handles the request for deleting an Account in the inventory
//
//	@Summary		Deletes an Account in the inventory
//	@Description	Deletes an Account present in the inventory by its ID
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			account_id	path		string	true	"Account ID"
//	@Success		200			{object}	nil
//	@Failure		500			{object}	GenericErrorResponse
//	@Router			/accounts/{account_id} [delete]
func (a APIServer) HandlerDeleteAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	a.logger.Debug("Removing an Account", zap.String("account_id", accountID))

	if err := a.sql.DeleteAccount(accountID); err != nil {
		a.logger.Error("Can't delete Cluster from DB", zap.String("account_id", accountID), zap.Error(err))
		c.PureJSON(http.StatusInternalServerError, NewGenericErrorResponse(err.Error()))
		return
	}

	c.PureJSON(http.StatusOK, apiresponsetypes.DeleteResponse{
		Count:  1,
		Status: fmt.Sprintf("Account '%s' Delete OK", accountID),
	})
}

// HandlerPatchAccount handles the request for patching an Account in the inventory
//
//	@Summary		Patches an Account in the inventory
//	@Description	Receives and patch into the DB the information for an existing Account
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			Account		body		inventory.Account	true	"Account to be modified"
//	@Param			account_id	path		string				true	"Account Name"
//	@Failure		501			{object}	nil					"Not Implemented"
//	@Router			/accounts/{account_id} [patch]
func (a APIServer) HandlerPatchAccount(c *gin.Context) {
	accountID := c.Param("account_id")
	a.logger.Debug("Patching an Account", zap.String("account", accountID))

	c.PureJSON(http.StatusNotImplemented, nil)
}
