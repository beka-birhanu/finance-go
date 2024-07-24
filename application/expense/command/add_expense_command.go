package command

import (
	"time"

	"github.com/google/uuid"
)

// AddExpenseCommand represents the command to add an expense.
type AddExpenseCommand struct {
	// UserId: The unique identifier of the user to whom the expense belongs.
	UserId uuid.UUID

	// Date: The date when the expense occurred.
	Date time.Time

	// Description: A brief description of the expense.
	Description string

	// Amount: The amount of the expense. Must be a positive float value.
	Amount float32
}
