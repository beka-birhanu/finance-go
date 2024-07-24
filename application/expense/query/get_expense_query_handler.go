package query

import (
	"github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/beka-birhanu/finance-go/domain/model"
)

type GetSingleExpenseQueryHandler struct {
	ExpenseRepository repository.IUserRepository
}

// Handle implements query.IQueryHandler.
func (h *GetSingleExpenseQueryHandler) Handle(query *GetExpenseQuery) (*model.Expense, error) {
	panic("unimplemented")
}

var _ query.IQueryHandler[*GetExpenseQuery, *model.Expense] = &GetSingleExpenseQueryHandler{}

func (h *GetExpenseQuery) Handle(query *GetExpenseQuery) (*model.Expense, error) {
	return nil, nil
}
