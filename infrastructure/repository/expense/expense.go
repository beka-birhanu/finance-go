package expenserepo

import (
	"database/sql"

	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	errexpense "github.com/beka-birhanu/finance-go/domain/error/expense"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

type PrimaryKey struct {
	Id     uuid.UUID
	UserId uuid.UUID
}

// used by both expense and user repo
var Expenses = map[PrimaryKey]expensemodel.Expense{}

var _ irepository.IExpenseRepository = &Repository{}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Save inserts or updates an expense in the repository.
// If the expense already exists, it updates the existing record.
// If the expense does not exist, it adds a new record.
//
// Returns:
//   - error: An error if any occurs, otherwise nil.
func (e *Repository) Save(expense *expensemodel.Expense) error {
	key := PrimaryKey{
		Id:     expense.ID(),
		UserId: expense.UserID(),
	}

	Expenses[key] = *expense
	return nil
}

// ById retrieves an expense by its unique identifier and user ID.
//
// Returns:
//   - *expensemodel.Expense: A pointer to the retrieved expense model.
//   - error: An error if the expense is not found, otherwise nil.
func (e *Repository) ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error) {
	key := PrimaryKey{
		Id:     id,
		UserId: userId,
	}

	expense, exists := Expenses[key]
	if !exists {
		return nil, errexpense.NotFound
	}

	return &expense, nil
}

