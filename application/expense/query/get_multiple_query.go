package expensqry

import (
	"time"

	"github.com/google/uuid"
)

// GetMultipleQuery represents a query for retrieving multiple expenses.
type GetMultipleQuery struct {
	UserID       uuid.UUID  // ID of the user whose expenses are to be retrieved
	Limit        int        // Maximum number of expenses to retrieve
	By           string     // Field to sort by (e.g., "date", "amount")
	LastSeenID   *uuid.UUID // ID of the last seen expense (for pagination)
	LastSeenTime *time.Time // Time of the last seen expense (for pagination)
	LastSeenAmt  float64    // Amount of the last seen expense (for pagination)
	Ascending    bool       // Whether to sort in ascending order
}

