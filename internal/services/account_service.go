package services

import (
	"context"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// AccountService defines the interface for account-related business logic.
type AccountService interface {
	List(ctx context.Context, options models.ListOptions) ([]db.AccountDBResponse, int, error)
	GetByID(ctx context.Context, accountID string) (db.AccountDBResponse, error)
	GetAccountClustersByID(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error)
	GetExpenseUpdateInstances(ctx context.Context, accountID string) ([]db.InstanceDBResponse, error)
	Create(ctx context.Context, accounts []inventory.Account) error
	Delete(ctx context.Context, accountID string) error
}

var _ AccountService = (*accountServiceImpl)(nil)

type accountServiceImpl struct {
	repo repositories.AccountRepository
	// other dependencies like other services or clients
}

// NewAccountService creates a new instance of AccountService.
func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountServiceImpl{
		repo: repo,
	}
}

// List retrieves a paginated list of accounts.
func (s *accountServiceImpl) List(ctx context.Context, options models.ListOptions) ([]db.AccountDBResponse, int, error) {
	return s.repo.ListAccounts(ctx, options)
}

// GetByID retrieves a single account by its name.
// It returns an error if no account or more than one account is found.
func (s *accountServiceImpl) GetByID(ctx context.Context, accountID string) (db.AccountDBResponse, error) {
	return s.repo.GetAccountByID(ctx, accountID)
}

// GetAccountClustersByID retrieves a single account by its name.
// It returns an error if no account or more than one account is found.
func (s *accountServiceImpl) GetAccountClustersByID(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error) {
	return s.repo.GetAccountClustersByID(ctx, accountID)
}

// GetExpenseUpdateInstances retrieves a single account by its name.
// It returns an error if no account or more than one account is found.
func (s *accountServiceImpl) GetExpenseUpdateInstances(ctx context.Context, accountID string) ([]db.InstanceDBResponse, error) {
	return s.repo.GetExpenseUpdateInstances(ctx, accountID)
}

// Create creates one or more new accounts.
func (s *accountServiceImpl) Create(ctx context.Context, accounts []inventory.Account) error {
	return s.repo.CreateAccount(ctx, accounts)
}

// Delete removes an account by its name.
func (s *accountServiceImpl) Delete(ctx context.Context, accountID string) error {
	return s.repo.DeleteAccount(ctx, accountID)
}
