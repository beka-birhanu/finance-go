package dto

import (
	"time"
)

type PatchRequest struct {
	Description *string    `json:"description,omitempty" validate:"omitempty"`
	Amount      *float32   `json:"amount,omitempty" validate:"omitempty"`
	Date        *time.Time `json:"date,omitempty" validate:"omitempty"`
}
