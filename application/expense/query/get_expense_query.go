package expense_queries

import "github.com/google/uuid"

type GetSingleExpenseQuery struct {
	UserId    uuid.UUID
	ExpenseId uuid.UUID
}
