package services

import (
	"context"
	"fmt"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/dto"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/repositories"
)

// AccountService defines the interface for account-related business logic.
type AccountService interface {
	List(ctx context.Context, options models.ListOptions) ([]db.AccountDBResponse, int, error)
	GetByID(ctx context.Context, accountID string) (db.AccountDBResponse, error)
	GetAccountClustersByID(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error)
	GetExpenseUpdateInstances(ctx context.Context, accountID string) ([]db.InstancePendingExpenseDB, error)
	Create(ctx context.Context, accounts []inventory.Account) error
	Update(ctx context.Context, accountID string, patch dto.AccountPatchRequest) error
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

// GetAccountClustersByID retrieves the clusters belonging to an account.
func (s *accountServiceImpl) GetAccountClustersByID(ctx context.Context, accountID string) ([]db.ClusterDBResponse, error) {
	clusters, err := s.repo.GetAccountClustersByID(ctx, accountID)
	if err != nil {
		return clusters, fmt.Errorf("get clusters for account %s: %w", accountID, err)
	}
	return clusters, nil
}

// GetExpenseUpdateInstances retrieves instances with outdated billing information.
func (s *accountServiceImpl) GetExpenseUpdateInstances(ctx context.Context, accountID string) ([]db.InstancePendingExpenseDB, error) {
	instances, err := s.repo.GetExpenseUpdateInstances(ctx, accountID)
	if err != nil {
		return instances, fmt.Errorf("get expense update instances for account %s: %w", accountID, err)
	}
	return instances, nil
}

// Create creates one or more new accounts.
func (s *accountServiceImpl) Create(ctx context.Context, accounts []inventory.Account) error {
	if err := s.repo.CreateAccount(ctx, accounts); err != nil {
		return fmt.Errorf("create accounts: %w", err)
	}
	return nil
}

// Update updates mutable fields of an existing account.
func (s *accountServiceImpl) Update(ctx context.Context, accountID string, patch dto.AccountPatchRequest) error {
	// Verify account exists before updating
	_, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	return s.repo.UpdateAccount(ctx, accountID, patch)
}

// Delete removes an account by its ID.
func (s *accountServiceImpl) Delete(ctx context.Context, accountID string) error {
	if err := s.repo.DeleteAccount(ctx, accountID); err != nil {
		return fmt.Errorf("delete account %s: %w", accountID, err)
	}
	return nil
}
