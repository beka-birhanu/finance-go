package query

import "github.com/google/uuid"

type GetSingleExpenseQuery struct {
	UserId    uuid.UUID
	ExpenseId uuid.UUID
}
