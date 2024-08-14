// Package expensqry provides functionality for handling queries related to expenses.
package expensqry

import (
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
)

// GetHandler processes queries to retrieve a specific expense.
type GetHandler struct {
	expenseRepository irepository.IExpenseRepository
}

// Ensure GetHandler implements iquery.IHandler interface for GetQuery.
var _ iquery.IHandler[*GetQuery, *expensemodel.Expense] = &GetHandler{}

// NewGetHandler creates a new instance of GetHandler with the provided expense repository.
func NewGetHandler(expenseRepository irepository.IExpenseRepository) *GetHandler {
	return &GetHandler{expenseRepository: expenseRepository}
}

// Handle retrieves an expense based on the provided query parameters.
//
// Returns:
//   - *expensemodel.Expense: The retrieved expense if found.
//   - error: An error if the retrieval fails.
func (h *GetHandler) Handle(query *GetQuery) (*expensemodel.Expense, error) {
	return h.expenseRepository.ById(query.ExpenseId, query.UserId)
}
