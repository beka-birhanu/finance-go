package dto

import (
	"time"
)

type AddExpenseRequest struct {
	UserId string    `json:"userId" validate:"required"`
	Title  string    `json:"title" validate:"required"`
	Amount float32   `json:"amount" validate:"required"`
	Data   time.Time `json:"date" validate:"required"`
}
