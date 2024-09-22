// Package irepository provides interfaces for accessing and managing expense data.
package irepository

import (
	"time"

	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
)

// ListByTimeParams defines parameters for retrieving expenses by time.
type ListByTimeParams struct {
	UserID       uuid.UUID  // ID of the user
	Limit        int        // Max number of expenses to return
	LastSeenID   *uuid.UUID // Pagination: ID of the last seen expense
	LastSeenDate *time.Time // Pagination: Time of the last seen expense
	Ascending    bool       // Sort order: true for ascending
}

// ListByAmountParams defines parameters for retrieving expenses by amount.
type ListByAmountParams struct {
	UserID      uuid.UUID  // ID of the user
	Limit       int        // Max number of expenses to return
	LastSeenID  *uuid.UUID // Pagination: ID of the last seen expense
	LastSeenAmt float64    // Pagination: Amount of the last seen expense
	Ascending   bool       // Sort order: true for ascending
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
