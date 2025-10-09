package handlers

import (
	"net/http"
	"strconv"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ExpenseHandler exposes expense-related HTTP endpoints.
type ExpenseHandler struct {
	service services.ExpenseService
	logger  *zap.Logger
}

// NewExpenseHandler constructs an ExpenseHandler.
func NewExpenseHandler(service services.ExpenseService, logger *zap.Logger) *ExpenseHandler {
	return &ExpenseHandler{
		service: service,
		logger:  logger,
	}
}

// expenseFilterParams defines the supported filter parameters
type expenseFilterParams struct {
	InstanceID string `form:"instance_id"`
}

// toRepoFilters converts bound query params into repository filters.
func (f *expenseFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.InstanceID != "" {
		filters["instance_id"] = f.InstanceID
	}
	return filters
}

// listExpensesRequest contains pagination and filters for List.
type listExpensesRequest struct {
	dto.PaginationRequest
	Filters expenseFilterParams `form:"inline"`
}

// List returns a paginated list of expenses.
//
//	@Summary		List expenses
//	@Description	Paginated retrieval with optional filters.
//	@Tags			Expenses
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"			default(1)
//	@Param			page_size	query		int		false	"Items per page"		default(10)
//	@Param			instance_id	query		string	false	"Instance ID filter"
//	@Success		200			{object}	responsetypes.ListResponse[dto.Expense]
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/expenses [get]
func (h *ExpenseHandler) List(c *gin.Context) {
	var req listExpensesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid query parameters: " + err.Error(),
		})
		return
	}

	opts := models.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	expenses, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing expenses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to list expenses",
		})
		return
	}

	response := responsetypes.NewListResponse(db.ToExpenseDTOResponseList(expenses), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Create inserts one or more expense records.
//
//	@Summary		Create expenses
//	@Description	Create one or multiple expense records.
//	@Tags			Expenses
//	@Accept			json
//	@Produce		json
//	@Param			expenses	body		[]dto.ExpenseDTORequest	true	"Expenses to create"
//	@Success		201			{object}	responsetypes.PostResponse
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/expenses [post]
func (h *ExpenseHandler) Create(c *gin.Context) {
	var expenseDTOs []dto.ExpenseDTORequest

	if err := c.ShouldBindJSON(&expenseDTOs); err != nil {
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), *dto.ToInventoryExpenseList(expenseDTOs)); err != nil {
		h.logger.Error("error creating expense", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create expenses: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(expenseDTOs),
		Status: "OK"},
	)
}
