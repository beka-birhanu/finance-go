// Package expensqry provides functionality for handling queries related to retrieving expenses.
package expensqry

import (
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

const (
	defaultLimit = 10       // Default limit for the number of expenses to retrieve
	minLimit     = 5        // Minimum limit for the number of expenses
	maxLimit     = 100      // Maximum limit for the number of expenses
	sortByAmount = "amount" // Field used for sorting by amount
)

// GetMultipleHandler handles queries for retrieving multiple expenses.
type GetMultipleHandler struct {
	expenseRepository irepository.IExpenseRepository // Repository for accessing expense data
}

// NewGetMultipleHandler creates a new instance of GetMultipleHandler with the given repository.
func NewGetMultipleHandler(expenseRepository irepository.IExpenseRepository) *GetMultipleHandler {
	return &GetMultipleHandler{expenseRepository: expenseRepository}
}

// Handle processes a GetMultipleQuery to retrieve multiple expenses based on the provided query parameters.
//
// Returns:
// - []*expensemodel.Expense: A slice of pointers to Expense models that match the query.
// - error: An error if the retrieval fails, such as issues with accessing the repository.
func (h *GetMultipleHandler) Handle(query *GetMultipleQuery) ([]*expensemodel.Expense, error) {
	// Set default limit if not provided
	limit := defaultLimit
	if query.Limit > 0 {
		if query.Limit < minLimit {
			limit = minLimit
		} else if query.Limit > maxLimit {
			limit = maxLimit
		} else {
			limit = query.Limit
		}
	}

	// Use amount-based pagination if specified
	if query.By == sortByAmount {
		return h.expenseRepository.ListByAmount(irepository.ListByAmountParams{
			UserID:      query.UserID,
			Limit:       limit,
			LastSeenID:  query.LastSeenID,
			LastSeenAmt: query.LastSeenAmt,
			Ascending:   query.Ascending,
		})
	}

	// Default to time-based pagination
	return h.expenseRepository.ListByTime(irepository.ListByTimeParams{
		UserID:       query.UserID,
		Limit:        limit,
		LastSeenID:   query.LastSeenID,
		LastSeenTime: query.LastSeenTime,
		Ascending:    query.Ascending,
	})
}
