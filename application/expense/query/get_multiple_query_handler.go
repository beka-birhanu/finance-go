package expensqry

import (
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

type GetMultipleHandler struct {
	expenseRepository irepository.IExpenseRepository
}

func NewGetMultipleHandler(expenseRepository irepository.IExpenseRepository) *GetMultipleHandler {
	return &GetMultipleHandler{expenseRepository: expenseRepository}
}

func (h *GetMultipleHandler) Handle(getMultipleuery *GetMultipleQuery) ([]*expensemodel.Expense, error) {
	return h.expenseRepository.List(getMultipleuery.UserId)
}
