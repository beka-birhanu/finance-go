package expensqry

import (
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

const (
	defaultLimit = 10
	minLimit     = 5
	maxLimit     = 100
	sortByAmount = "amount"
)

type GetMultipleHandler struct {
	expenseRepository irepository.IExpenseRepository
}

func NewGetMultipleHandler(expenseRepository irepository.IExpenseRepository) *GetMultipleHandler {
	return &GetMultipleHandler{expenseRepository: expenseRepository}
}

func (h *GetMultipleHandler) Handle(query *GetMultipleQuery) ([]*expensemodel.Expense, error) {
	// Set defaults if not provided
	limit := defaultLimit
	if query.Limit > 0 {
		// Apply max and min limits
		if query.Limit < minLimit {
			limit = minLimit
		} else if query.Limit > maxLimit {
			limit = maxLimit
		} else {
			limit = query.Limit
		}
	}

	// Default to time-based pagination if By is not sortByAmount
	if query.By == sortByAmount {
		// Handle amount-based pagination
		return h.expenseRepository.ListByAmount(irepository.ListByAmountParams{
			UserID:      query.UserID,
			Limit:       limit,
			LastSeenID:  query.LastSeenID,
			LastSeenAmt: query.LastSeenAmt,
			Ascending:   query.Ascending,
		})
	}

	// Handle time-based pagination as default
	return h.expenseRepository.ListByTime(irepository.ListByTimeParams{
		UserID:       query.UserID,
		Limit:        limit,
		LastSeenID:   query.LastSeenID,
		LastSeenTime: query.LastSeenTime,
		Ascending:    query.Ascending,
	})
}
