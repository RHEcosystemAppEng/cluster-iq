package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
