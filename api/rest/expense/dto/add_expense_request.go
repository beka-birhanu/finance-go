package dto

import (
	"time"
)

type AddExpenseRequest struct {
	Description string    `json:"description" validate:"required"`
	Amount      float32   `json:"amount" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}
