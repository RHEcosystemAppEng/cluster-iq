package repositories

import (
	"context"
	"fmt"

	dbclient "github.com/RHEcosystemAppEng/cluster-iq/internal/db_client"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models"
	"github.com/RHEcosystemAppEng/cluster-iq/internal/models/db"
)

const (
	// DB Table for expenses
	ExpensesTable = "expenses"
	// InsertExpensesQuery inserts into a new expense for an instance
	InsertExpensesQuery = `
		INSERT INTO expenses (
			instance_id,
			date,
			amount
		) VALUES (
			(SELECT id FROM instances WHERE instance_id=:instance_id),
			:date,
			:amount
		) ON CONFLICT (instance_id, date) DO UPDATE SET
			amount = EXCLUDED.amount
	`
)

var _ ExpenseRepository = (*expenseRepositoryImpl)(nil)

// ExpenseRepository defines the interface for data access operations for expenses.
type ExpenseRepository interface {
	ListExpenses(ctx context.Context, opts models.ListOptions) ([]db.ExpenseDBResponse, int, error)
	GetExpensesByInstance(ctx context.Context, instanceID string) ([]db.ExpenseDBResponse, error)
	Create(ctx context.Context, expenses []inventory.Expense) error
}

type expenseRepositoryImpl struct {
	db *dbclient.DBClient
}

func NewExpenseRepository(db *dbclient.DBClient) ExpenseRepository {
	return &expenseRepositoryImpl{db: db}
}

// ListExpenses retrieves all expenses from the database.
//
// Parameters:
//
// Returns:
// - A slice of inventory.Expense objects.
// - An error if the query fails.
func (r *expenseRepositoryImpl) ListExpenses(ctx context.Context, opts models.ListOptions) ([]db.ExpenseDBResponse, int, error) {
	var expenses []db.ExpenseDBResponse

	if err := r.db.SelectWithContext(ctx, &expenses, ExpensesTable, opts, "date", "*"); err != nil {
		return expenses, 0, fmt.Errorf("failed to list expenses: %w", err)
	}

	return expenses, len(expenses), nil
}

// GetExpensesByInstance retrieves expenses for a specific instance.
//
// Parameters:
// - instanceID: The ID of the instance.
//
// Returns:
// - A slice of inventory.Expense objects associated with the instance.
// - An error if the query fails.
func (r *expenseRepositoryImpl) GetExpensesByInstance(ctx context.Context, instanceID string) ([]db.ExpenseDBResponse, error) {
	opts := models.ListOptions{
		PageSize: 0,
		Offset:   0,
		Filters: map[string]interface{}{
			"instance_id": instanceID,
		},
	}

	expenses, _, err := r.ListExpenses(ctx, opts)
	return expenses, err
}

// Create writes a batch of expenses to the database in a transaction.
//
// Parameters:
// - expenses: A slice of inventory.Expense objects to insert.
//
// Returns:
// - An error if the transaction fails.
func (r *expenseRepositoryImpl) Create(ctx context.Context, expenses []inventory.Expense) error {
	if err := r.db.InsertWithContext(ctx, InsertExpensesQuery, expenses); err != nil {
		return err
	}
	return nil
}
