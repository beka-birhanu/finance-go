package expensqry

import "github.com/google/uuid"

type GetQuery struct {
	UserId    uuid.UUID
	ExpenseId uuid.UUID
}
