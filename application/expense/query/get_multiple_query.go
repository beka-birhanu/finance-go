package expensqry

import (
	"time"

	"github.com/google/uuid"
)

type GetMultipleQuery struct {
	UserID       uuid.UUID
	Limit        int
	By           string
	LastSeenID   *uuid.UUID
	LastSeenTime *time.Time
	LastSeenAmt  float64
	Ascending    bool
}
