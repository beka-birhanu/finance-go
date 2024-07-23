package expense_queries

import (
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/models"
)

type GetSingleExpenseQueryHandler struct {
	ExpenseRepository persistance.IUserRepository
}

func (h *GetSingleExpenseQuery) Handle(query *GetSingleExpenseQuery) (*models.Expense, error) {
	return nil, nil
}
