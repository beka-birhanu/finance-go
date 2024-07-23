package query

import (
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/beka-birhanu/finance-go/domain/model"
)

type GetSingleExpenseQueryHandler struct {
	ExpenseRepository repository.IUserRepository
}

func (h *GetSingleExpenseQuery) Handle(query *GetSingleExpenseQuery) (*model.Expense, error) {
	return nil, nil
}
