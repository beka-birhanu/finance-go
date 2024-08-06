package expensqry

import "github.com/google/uuid"

// GetQuery represents a query for retrieving a specific expense.
type GetQuery struct {
	UserId    uuid.UUID // ID of the user
	ExpenseId uuid.UUID // ID of the expense
}
