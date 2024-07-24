/*
Package expensemodel includes the definition of the Expense aggregate, which represents
an individual expense, and provides functions for creating and interacting with expenses.

Key Components:
- Expense: Represents an expense with details such as description, amount, and
associated user.
- Config: Holds the mandatory parameters required to create a new Expense.
- New: Creates a new Expense instance based on the provided configuration.
- validateDescription: Validates that the expense description meets length constraints.

Dependencies:
- github.com/google/uuid: Used for generating unique IDs.
- time: Used for timestamps.
*/
package expensemodel

import (
	"strings"
	"time"

	"github.com/beka-birhanu/finance-go/domain/error/expense"
	"github.com/google/uuid"
)

const (
	maxDescriptionLength = 255
)

// Expense represents an expense aggregate.
type Expense struct {
	id          uuid.UUID
	description string
	amount      float32
	date        time.Time
	userId      uuid.UUID
	createdAt   time.Time
	updatedAt   time.Time
}

// Config holds all mandatory parameters for creating a new Expense.
type Config struct {
	// Description must be non-empty and adhere to length constraints.
	Description string

	// Amount must be a positive number.
	Amount float32

	// UserId is the ID of the owner user for the expense.
	UserId uuid.UUID

	// Date is the timestamp when the expense occurred.
	Date time.Time

	// CreationTime is the timestamp when the expense is created.
	CreationTime time.Time
}

// New creates a new Expense with the provided configuration.
//
// Returns:
// - A pointer to the newly created Expense if successful.
// - An error if any of the following conditions are not met:
//   - Any field in the config is missing or invalid.
//   - The description does not meet length constraints.
//   - The amount is not positive.
func New(config Config) (*Expense, error) {
	config.Description = strings.TrimSpace(config.Description)
	if err := validateDescription(config.Description); err != nil {
		return nil, err
	}

	if config.Amount <= 0 {
		return nil, errexpense.NegativeAmount
	}

	return &Expense{
		id:          uuid.New(), // New ID for the expense
		description: config.Description,
		amount:      config.Amount,
		userId:      config.UserId,
		date:        config.Date,
		createdAt:   config.CreationTime,
		updatedAt:   config.CreationTime,
	}, nil
}

func validateDescription(desc string) error {
	if desc == "" {
		return errexpense.EmptyDescription
	}

	if len(desc) > maxDescriptionLength {
		return errexpense.DescriptionTooLong
	}

	return nil
}

// ID returns the ID of the expense.
func (e *Expense) ID() uuid.UUID {
	return e.id
}

// Description returns the description of the expense.
func (e *Expense) Description() string {
	return e.description
}

// Amount returns the amount of the expense.
func (e *Expense) Amount() float32 {
	return e.amount
}

// Date returns the date of the expense.
func (e *Expense) Date() time.Time {
	return e.date
}

// UserID returns the ID of the user associated with the expense.
func (e *Expense) UserID() uuid.UUID {
	return e.userId
}

// CreatedAt returns the creation timestamp of the expense.
func (e *Expense) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returns the last update timestamp of the expense.
func (e *Expense) UpdatedAt() time.Time {
	return e.updatedAt
}

