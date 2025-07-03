package handlers

import (
	"net/http"
	"strconv"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/mappers"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
)

// ExpenseHandler handles HTTP requests for expenses.
type ExpenseHandler struct {
	service services.ExpenseService
}

func NewExpenseHandler(service services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
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

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	expenses, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list expenses"))
		return
	}

	expenseDTOs := mappers.ToExpenseDTOs(expenses)
	response := dto.NewListResponse(expenseDTOs, total)
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}
