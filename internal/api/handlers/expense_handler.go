package handlers

import (
	"net/http"

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

type listExpensesRequest struct {
	dto.PaginationRequest
	InstanceID string `form:"instance_id"`
}

// List handles the request to list all expenses.
//
//	@Summary		List all expenses
//	@Description	Returns a paginated list of expenses
//	@Tags			Expenses
//	@Param			page		query	int		false	"Page number for pagination"	default(1)
//	@Param			pageSize	query	int		false	"Number of items per page"		default(10)
//	@Param			instance_id	query	string	false	"Filter by instance ID"
//	@Success		200			{object}	dto.ListResponse[dto.Expense]
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/expenses [get]
func (h *ExpenseHandler) List(c *gin.Context) {
	var req listExpensesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	filters := make(map[string]interface{})
	if req.InstanceID != "" {
		filters["instance_id"] = req.InstanceID
	}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  filters,
	}

	expenses, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to list expenses"))
		return
	}

	expenseDTOs := mappers.ToExpenseDTOs(expenses)
	response := dto.NewListResponse(expenseDTOs, total)
	c.JSON(http.StatusOK, response)
}
