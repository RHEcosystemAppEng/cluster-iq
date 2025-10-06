package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/api/mappers"
	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AccountHandler handles HTTP requests for accounts.
type AccountHandler struct {
	service services.AccountService
	logger  *zap.Logger
}

func NewAccountHandler(service services.AccountService, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{
		service: service,
		logger:  logger,
	}
}

type accountFilterParams struct {
	Provider string `form:"provider"`
}

func (f *accountFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.Provider != "" {
		filters["provider"] = f.Provider
	}
	return filters
}

type listAccountsRequest struct {
	dto.PaginationRequest
	Filters accountFilterParams `form:"inline"`
}

// List handles the request for obtaining the Account list.
//
//	@Summary		List accounts
//	@Description	Returns a paginated list of accounts based on optional filters.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number for pagination"	default(1)
//	@Param			page_size	query		int		false	"Number of items per page"		default(10)
//	@Param			provider	query		string	false	"Filter by cloud provider"		example(aws)
//	@Success		200			{object}	dto.ListResponse[dto.Account]
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
	var req listAccountsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid query parameters: "+err.Error()))
		return
	}

	opts := models.ListOptions{
		PageSize: req.PageSize,
		Offset:   (req.Page - 1) * req.PageSize,
		Filters:  req.Filters.toRepoFilters(),
	}

	accounts, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve accounts"))
		return
	}

	response := dto.NewListResponse(mappers.ToAccountDTOResponseList(accounts), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// GetByID handles the request for obtaining a single account by its ID.
//
//	@Summary		Get an account by name
//	@Description	Returns a single account.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"Account Name"
//	@Success		200		{object}	dto.Account
//	@Failure		404		{object}	dto.GenericErrorResponse
//	@Failure		500		{object}	dto.GenericErrorResponse
//	@Router			/accounts/{name} [get]
func (h *AccountHandler) GetByID(c *gin.Context) {
	accountID := c.Param("id")

	account, err := h.service.GetByID(c.Request.Context(), accountID)
	if err != nil {
		h.logger.Error("error getting an account", zap.String("account_id", accountID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Account not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve account"))
		return
	}

	c.JSON(http.StatusOK, account.ToAccountDTOResponse())
}

// GetAccountClustersByID handles the request for obtaining a single account by its ID.
//
//	@Summary		Get an account by name
//	@Description	Returns a single account.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string	true	"Account Name"
//	@Success		200		{object}	dto.Account
//	@Failure		404		{object}	dto.GenericErrorResponse
//	@Failure		500		{object}	dto.GenericErrorResponse
//	@Router			/accounts/{name} [get]
func (h *AccountHandler) GetAccountClustersByID(c *gin.Context) {
	accountID := c.Param("id")

	clusters, err := h.service.GetAccountClustersByID(c.Request.Context(), accountID)
	if err != nil {
		h.logger.Error("error getting an account", zap.String("account_id", accountID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Account not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to retrieve account"))
		return
	}

	response := dto.NewListResponse(mappers.ToClusterDTOResponseList(clusters), len(clusters))

	c.JSON(http.StatusOK, response)
}

// Create handles the creation of new accounts.
//
//	@Summary		Create accounts
//	@Description	Creates one or more new accounts.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			accounts	body		[]dto.NewAccount	true	"Account or accounts to create"
//	@Success		201			{object}	[]dto.Account
//	@Failure		400			{object}	dto.GenericErrorResponse
//	@Failure		500			{object}	dto.GenericErrorResponse
//	@Router			/accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
	var newAccountsDTO []dto.AccountDTORequest
	if err := c.ShouldBindJSON(&newAccountsDTO); err != nil {
		h.logger.Error("error processing received accounts", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.NewGenericErrorResponse("Invalid request body: "+err.Error()))
		return
	}

	if err := h.service.Create(c.Request.Context(), mappers.ToAccountModelList(newAccountsDTO)); err != nil {
		h.logger.Error("error creating accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to create accounts: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(newAccountsDTO),
		Status: "OK"},
	)
}

// Delete handles the deletion of an account.
//
//	@Summary		Delete an account
//	@Description	Deletes an account by its name.
//	@Tags			Accounts
//	@Accept			json
//	@Param			name	path		string	true	"Account Name"
//	@Success		204		{object}	nil
//	@Failure		404		{object}	dto.GenericErrorResponse
//	@Failure		500		{object}	dto.GenericErrorResponse
//	@Router			/accounts/{name} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	accountID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), accountID); err != nil {
		h.logger.Error("error deleting account", zap.String("account_id", accountID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.NewGenericErrorResponse("Account not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.NewGenericErrorResponse("Failed to delete account: "+err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

// Update handles the update of an existing account.
//
//	@Summary		Update an account
//	@Description	Updates an existing account by its name.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string		true	"Account Name"
//	@Param			account	body		dto.Account	true	"Updated account data"
//	@Success		200		{object}	nil
//	@Failure		501		{object}	nil	"Not Implemented"
//	@Router			/accounts/{name} [patch]
func (h *AccountHandler) Update(c *gin.Context) {
	// TODO
	c.PureJSON(http.StatusNotImplemented, nil)
}
