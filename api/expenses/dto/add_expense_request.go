package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddExpenseRequest struct {
	UserId      uuid.UUID `json:"userId" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Amount      float32   `json:"amount" validate:"required"`
	Data        time.Time `json:"date" validate:"required"`
}
