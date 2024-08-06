package expensecmd

import (
	"time"

	"github.com/google/uuid"
)

// PatchCommand represents a command to update an existing expense.
type PatchCommand struct {
	Description *string    // Optional new description for the expense
	Amount      *float32   // Optional new amount for the expense
	Date        *time.Time // Optional new date for the expense
	Id          uuid.UUID  // Unique identifier of the expense to be updated
	UserId      uuid.UUID  // Identifier of the user who owns the expense
}
