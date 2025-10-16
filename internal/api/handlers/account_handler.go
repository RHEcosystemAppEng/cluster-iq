package handlers

import (
	"errors"
	"net/http"
	"strconv"

	responsetypes "github.com/RHEcosystemAppEng/cluster-iq/internal/api/response_types"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AccountHandler connects HTTP endpoints with the AccountService.
type AccountHandler struct {
	service services.AccountService
	logger  *zap.Logger
}

// NewAccountHandler builds a new AccountHandler with its dependencies.
func NewAccountHandler(service services.AccountService, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{
		service: service,
		logger:  logger,
	}
}

// accountFilterParams defines the supported filter parameters
type accountFilterParams struct {
	Provider string `form:"provider"`
}

// toRepoFilters translates bound query filters into repository filter keys.
func (f *accountFilterParams) toRepoFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if f.Provider != "" {
		filters["provider"] = f.Provider
	}
	return filters
}

// listAccountsRequest binds pagination and filter query params for List.
type listAccountsRequest struct {
	dto.PaginationRequest
	Filters accountFilterParams `form:"inline"`
}

// List returns a paginated list of accounts.
//
//	@Summary		List accounts
//	@Description	Paginated retrieval of accounts with optional filters.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"		default(1)
//	@Param			page_size	query		int		false	"Items per page"	default(10)
//	@Param			provider	query		string	false	"Cloud provider"	example(aws)
//	@Success		200			{object}	responsetypes.ListResponse[dto.AccountDTOResponse]
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/accounts [get]
func (h *AccountHandler) List(c *gin.Context) {
	var req listAccountsRequest

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

	accounts, total, err := h.service.List(c.Request.Context(), opts)
	if err != nil {
		h.logger.Error("error listing accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve accounts",
		})
		return
	}

	response := responsetypes.NewListResponse(db.ToAccountDTOResponseList(accounts), total)

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, response)
}

// GetByID retrieves a single account by ID.
//
//	@Summary		Get an account by ID
//	@Description	Return a single account resource.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Account ID"
//	@Success		200	{object}	dto.AccountDTOResponse
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/accounts/{id} [get]
//
// NOTE: The handler reads path param "id". Ensure routing path uses {id}.
func (h *AccountHandler) GetByID(c *gin.Context) {
	accountID := c.Param("id")

	account, err := h.service.GetByID(c.Request.Context(), accountID)
	if err != nil {
		h.logger.Error("error getting an account", zap.String("account_id", accountID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve account",
		})
		return
	}

	c.JSON(http.StatusOK, account.ToAccountDTOResponse())
}

// GetAccountClustersByID lists clusters belonging to a given account ID.
//
//	@Summary		List clusters by account ID
//	@Description	Return the clusters associated with the specified account.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Account ID"
//	@Success		200	{object}	responsetypes.ListResponse[dto.ClusterDTOResponse]
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/accounts/{id}/clusters [get]
//
// NOTE: Align the documented route with the actual router configuration.
func (h *AccountHandler) GetAccountClustersByID(c *gin.Context) {
	accountID := c.Param("id")

	clusters, err := h.service.GetAccountClustersByID(c.Request.Context(), accountID)
	if err != nil {
		h.logger.Error("error getting an account", zap.String("account_id", accountID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to retrieve account",
		})
		return
	}

	response := responsetypes.NewListResponse(db.ToClusterDTOResponseList(clusters), len(clusters))

	c.JSON(http.StatusOK, response)
}

// Create creates one or more accounts.
//
//	@Summary		Create accounts
//	@Description	Create one or multiple accounts from the request body.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			accounts	body		[]dto.AccountDTORequest	true	"Accounts to create"
//	@Success		201			{object}	responsetypes.PostResponse
//	@Failure		400			{object}	responsetypes.GenericErrorResponse
//	@Failure		500			{object}	responsetypes.GenericErrorResponse
//	@Router			/accounts [post]
func (h *AccountHandler) Create(c *gin.Context) {
	var newAccountsDTO []dto.AccountDTORequest

	if err := c.ShouldBindJSON(&newAccountsDTO); err != nil {
		h.logger.Error("error processing received accounts", zap.Error(err))
		c.JSON(http.StatusBadRequest, responsetypes.GenericErrorResponse{
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), *dto.ToInventoryAccountList(newAccountsDTO)); err != nil {
		h.logger.Error("error creating accounts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to create accounts: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responsetypes.PostResponse{
		Count:  len(newAccountsDTO),
		Status: "OK"},
	)
}

// Delete removes an account by ID.
//
//	@Summary		Delete an account
//	@Description	Delete an account resource by its ID.
//	@Tags			Accounts
//	@Accept			json
//	@Param			id	path		string	true	"Account ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	responsetypes.GenericErrorResponse
//	@Failure		500	{object}	responsetypes.GenericErrorResponse
//	@Router			/accounts/{id} [delete]
func (h *AccountHandler) Delete(c *gin.Context) {
	accountID := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), accountID); err != nil {
		h.logger.Error("error deleting account", zap.String("account_id", accountID), zap.Error(err))
		if errors.Is(err, repositories.ErrNotFound) {
			c.JSON(http.StatusNotFound, responsetypes.GenericErrorResponse{
				Message: "Account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responsetypes.GenericErrorResponse{
			Message: "Failed to delete account: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Update applies partial updates to an existing account.
//
//	@Summary		Update an account
//	@Description	Patch an existing account by ID.
//	@Tags			Accounts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Account ID"
//	@Param			account	body		dto.AccountDTORequest	true	"Partial account payload"
//	@Success		200		{object}	nil
//	@Failure		501		{object}	nil	"Not Implemented"
//	@Router			/accounts/{id} [patch]
func (h *AccountHandler) Update(c *gin.Context) {
	// TODO: Implement partial update semantics
	c.PureJSON(http.StatusNotImplemented, nil)
}
