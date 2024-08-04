package irepository

import (
	"time"

	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

// ListByTimeParams holds parameters for ListByTime function
type ListByTimeParams struct {
	UserID       uuid.UUID
	Limit        int
	LastSeenID   *uuid.UUID
	LastSeenTime *time.Time
	Ascending    bool
}

// ListByAmountParams holds parameters for ListByAmount function
type ListByAmountParams struct {
	UserID      uuid.UUID
	Limit       int
	LastSeenID  *uuid.UUID
	LastSeenAmt float64
	Ascending   bool
}

// IExpenseRepository defines methods for accessing and managing expense data.
type IExpenseRepository interface {
	// Save inserts or updates an expense in the repository.
	Save(expense *expensemodel.Expense) error

	// ById retrieves an expense by its unique identifier and user ID.
	ById(id uuid.UUID, userId uuid.UUID) (*expensemodel.Expense, error)

	// ListByTime retrieves paginated expenses by user ID based on created time.
	ListByTime(params ListByTimeParams) ([]*expensemodel.Expense, error)

	// ListByAmount retrieves paginated expenses by user ID based on amount.
	ListByAmount(params ListByAmountParams) ([]*expensemodel.Expense, error)
}
