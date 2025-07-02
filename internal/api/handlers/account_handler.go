package handlers

import (
	"errors"
	"net/http"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/mappers"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
)

// AccountHandler handles HTTP requests for accounts.
type AccountHandler struct {
	service services.AccountService
}

func NewAccountHandler(service services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

type listAccountsRequest struct {
	dto.PaginationRequest
	Provider string `form:"provider"`
	Name     string `form:"name"`
	ID       string `form:"id"`
}

// List handles the request for obtaining the Account list.
//
//	@Summary		List accounts
//	@Description	Returns a paginated list of accounts based on optional filters.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Page number for pagination"	default(1)
//	@Param			pageSize	query	int		false	"Number of items per page"		default(10)
//	@Param			provider	query	string	false	"Filter by cloud provider"		example(aws)
//	@Param			name		query	string	false	"Filter by account name"
//	@Param			id			query	string	false	"Filter by account ID"
//	@Success		200			{object}	dto.ListResponse[dto.Account]
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
	var req listAccountsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	filters := make(map[string]interface{})
	if req.Provider != "" {
		filters["provider"] = req.Provider
	}
	if req.Name != "" {
		filters["name"] = req.Name
	}
	if req.ID != "" {
		filters["id"] = req.ID
	}

	opts := repositories.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  filters,
	}

	accounts, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve accounts"))
		return
	}

	accountDTOs := mappers.ToAccountDTOs(accounts)
	response := dto.NewListResponse(accountDTOs, total)
	c.JSON(http.StatusOK, response)
}

// GetByName handles the request for obtaining a single account by its name.
//
//	@Summary		Get an account by name
//	@Description	Returns a single account.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"Account Name"
//	@Success		200		{object}	dto.Account
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/accounts/{name} [get]
func (h *AccountHandler) GetByName(c *gin.Context) {
	accountName := c.Param("name")

	account, err := h.service.GetByName(c.Request.Context(), accountName)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) { // Assuming ErrNotFound exists
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Account not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve account"))
		return
	}

	accountDTO := mappers.ToAccountDTO(account)
	c.JSON(http.StatusOK, accountDTO)
}

// Create handles the creation of new accounts.
//
//	@Summary		Create accounts
//	@Description	Creates one or more new accounts.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			accounts	body		[]dto.NewAccount	true	"Account or accounts to create"
//	@Success		201			{object}	nil
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
	var newAccountsDTO []dto.NewAccount
	if err := c.ShouldBindJSON(&newAccountsDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	accounts := mappers.ToAccountModels(newAccountsDTO)
	if err := h.service.Create(c.Request.Context(), accounts); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create accounts: "+err.Error()))
		return
	}

	c.Status(http.StatusCreated)
}

// Delete handles the deletion of an account.
//
//	@Summary		Delete an account
//	@Description	Deletes an account by its name.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"Account Name"
//	@Success		204		{object}	nil
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/accounts/{name} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	accountName := c.Param("name")

	if err := h.service.Delete(c.Request.Context(), accountName); err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to delete account: "+err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}
