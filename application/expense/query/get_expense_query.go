package query

import "github.com/google/uuid"

type GetExpenseQuery struct {
	UserId    uuid.UUID
	ExpenseId uuid.UUID
}
