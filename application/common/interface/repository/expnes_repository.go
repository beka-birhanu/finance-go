package irepository

import (
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

// IExpenseRepository defines methods for accessing and managing expense data.
//
// Methods:
// - Add(expense *expensemodel.Expense) error: Creates a new expense in the repository and
// returns an error if any occurs.
// - ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error): Retrieves an expense by its unique
// ID and returns the expense or an error if the expense does not exist.
// - Update(expense *expensemodel.Expense) error: Updates the expense and returns an error if the expense does not exist.
type IExpenseRepository interface {
	// Add adds a new expense to the repository.
	Add(expense *expensemodel.Expense) error

	// Update updates the expense and returns an error if the expense does not exist.
	Update(expense *expensemodel.Expense) error

	// ById retrieves an expense by its unique identifier.
	ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error)
}
