// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateExpenseInput struct {
	Description string    `json:"description"`
	Amount      float32   `json:"amount"`
	Date        time.Time `json:"date"`
	UserID      uuid.UUID `json:"userId"`
}

type Expense struct {
	ID          uuid.UUID  `json:"id"`
	Description string     `json:"description"`
	Amount      float32    `json:"amount"`
	Date        time.Time  `json:"date"`
	UserID      uuid.UUID  `json:"userId"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type UpdateExpenseInput struct {
	Description *string    `json:"description,omitempty"`
	Amount      *float32   `json:"amount,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	UserID      uuid.UUID  `json:"userId"`
}
