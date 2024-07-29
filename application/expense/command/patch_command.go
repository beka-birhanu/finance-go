package expensecmd

import (
	"time"

	"github.com/google/uuid"
)

type PatchCommand struct {
	Description *string
	Amount      *float32
	Date        *time.Time
	Id          uuid.UUID
	UserId      uuid.UUID
}
