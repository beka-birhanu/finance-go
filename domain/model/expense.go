package model

import (
	"strings"
	"time"

	domainError "github.com/beka-birhanu/finance-go/domain/error"
	"github.com/google/uuid"
)

const (
	MAX_DESCRIPTION_LENGTH = 255
)

type Expense struct {
	id          uuid.UUID
	description string
	amount      float32
	date        time.Time
	userId      uuid.UUID
	createdAt   time.Time
	updatedAt   time.Time
}

func NewExpense(description string, amount float32, userId uuid.UUID, date time.Time) (*Expense, error) {
	description = strings.TrimSpace(description)

	if amount <= 0 {
		return nil, domainError.ErrNegativeExpenseAmount
	}

	if len(description) > MAX_DESCRIPTION_LENGTH {
		return nil, domainError.ErrExpenseDescriptionTooLong
	}

	if description == "" {
		return nil, domainError.ErrEmptyExpenseDescription
	}

	return &Expense{
		id:          uuid.New(),
		description: description,
		amount:      amount,
		userId:      userId,
		date:        date,
		createdAt:   time.Now().UTC(),
		updatedAt:   time.Now().UTC(),
	}, nil
}

func (e *Expense) ID() uuid.UUID {
	return e.id
}

func (e *Expense) Description() string {
	return e.description
}

func (e *Expense) Amount() float32 {
	return e.amount
}

func (e *Expense) Date() time.Time {
	return e.date
}

func (e *Expense) UserID() uuid.UUID {
	return e.userId
}

func (e *Expense) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Expense) UpdatedAt() time.Time {
	return e.updatedAt
}
