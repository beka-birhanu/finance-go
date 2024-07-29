package irepository

import (
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

// IExpenseRepository defines methods for accessing and managing expense data.
//
// Methods:
// - Save(expense *expensemodel.Expense) error: Inserts or updates an expense in the repository and
// returns an error if any occurs.
// - ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error): Retrieves an expense by its unique
// ID and returns the expense or an error if the expense does not exist.
type IExpenseRepository interface {
	// Save inserts or updates an expense in the repository.
	Save(expense *expensemodel.Expense) error

	// ById retrieves an expense by its unique identifier and user ID.
	ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error)

	// List retrieves expenses by user ID.
	List(userId uuid.UUID) ([]*expensemodel.Expense, error)
}
