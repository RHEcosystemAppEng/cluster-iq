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

// ExpenseHandler handles HTTP requests for expenses.
type ExpenseHandler struct {
	service services.ExpenseService
	logger  *zap.Logger
}

func NewExpenseHandler(service services.ExpenseService, logger *zap.Logger) *ExpenseHandler {
	return &ExpenseHandler{
		service: service,
		logger:  logger,
	}
}

type expenseFilterParams struct {
	InstanceID string `form:"instance_id"`
}

func (f *expenseFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.InstanceID != "" {
		filters["instance_id"] = f.InstanceID
	}
	return filters
}

type listExpensesRequest struct {
	dto.PaginationRequest
	Filters expenseFilterParams `form:"inline"`
}

// List handles the request to list all expenses.
//
//	@Summary		List all expenses
//	@Description	Returns a paginated list of expenses
//	@Tags			Expenses
//	@Param			page		query		int		false	"Page number for pagination"	default(1)
//	@Param			page_size	query		int		false	"Number of items per page"		default(10)
//	@Param			instance_id	query		string	false	"Filter by instance ID"
//	@Success		200			{object}	dto.ListResponse[dto.Expense]
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/expenses [get]
func (h *ExpenseHandler) List(c *gin.Context) {
	var req listExpensesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
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
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list expenses"))
		return
	}

	response := dto.NewListResponse(db.ToExpenseDTOResponseList(expenses), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// Create handles the request to create new expense records.
//
//	@Summary		Create new expense records
//	@Description	Adds one or more expense records to the database
//	@Tags			Expenses
//	@Accept			json
//	@Produce		json
//	@Param			expenses	body		[]dto.CreateExpense	true	"A list of new expenses to create"
//	@Success		201			{object}	nil
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/expenses [post]
func (h *ExpenseHandler) Create(c *gin.Context) {
	var expenseDTOs []dto.ExpenseDTORequest
	if err := c.ShouldBindJSON(&expenseDTOs); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	if err := h.service.Create(c.Request.Context(), *dto.ToInventoryExpenseList(expenseDTOs)); err != nil {
		h.logger.Error("error creating expense", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create expenses: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(expenseDTOs),
		Status: "OK"},
	)
}
