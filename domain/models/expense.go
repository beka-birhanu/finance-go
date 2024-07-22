package models

import (
	"strings"
	"time"

	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/google/uuid"
)

const (
	MAX_DESCRIPTION_LENGTH = 255
)

type Expense struct {
	ID          uuid.UUID
	Description string
	Amount      float32
	Date        time.Time
	UserId      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewExpense(description string, amount float32, userId uuid.UUID) (*Expense, error) {
	description = strings.TrimSpace(description)

	if amount <= 0 {
		return nil, domain_errors.ErrNegativeExpenseAmount
	}

	if len(description) > MAX_DESCRIPTION_LENGTH {
		return nil, domain_errors.ErrExpenseDescriptionTooLong
	}

	if description == "" {
		return nil, domain_errors.ErrEmptyExpenseDescription
	}

	return &Expense{
		ID:          uuid.New(),
		Description: description,
		Amount:      amount,
		UserId:      userId,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

