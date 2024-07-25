package expensqry

import (
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

type GetHandler struct {
	ExpenseRepository irepository.IUserRepository
}

var _ iquery.IHandler[*GetQuery, *expensemodel.Expense] = &GetHandler{}

func NewGetHandler(expenseRepository irepository.IUserRepository) {
}
func (h *GetHandler) Handle(query *GetQuery) (*expensemodel.Expense, error) {
	return nil, nil
}
