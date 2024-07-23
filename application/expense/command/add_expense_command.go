package commands

import (
	"time"

	"github.com/google/uuid"
)

type AddExpenseCommand struct {
	UserId      uuid.UUID
	Date        time.Time
	Description string
	Amount      float32
}
