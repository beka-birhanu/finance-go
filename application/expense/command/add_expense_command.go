package expensecmd

import (
	"time"

	"github.com/google/uuid"
)

// AddCommand represents the command to add an expense.
type AddCommand struct {
	// UserId: The unique identifier of the user to whom the expense belongs.
	UserId uuid.UUID

	// Date: The date when the expense occurred.
	Date time.Time

	// Description: A brief description of the expense.
	Description string

	// Amount: The amount of the expense. Must be a positive float value.
	Amount float32
}
