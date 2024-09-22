package utils

import (
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/graph/model"
	"github.com/beka-birhanu/finance-go/api/utils"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// NewGQLError creates a new gqlerror with a custom message and extensions
func NewGQLError(err errapi.Error) *gqlerror.Error {
	var errMessage string
	extensions := make(map[string]interface{})
	extensions["StatusCode"] = err.StatusCode()

	switch err.StatusCode() {
	case errapi.BadRequest, errapi.Conflict, errapi.NotFound, errapi.Forbidden:
		errMessage = err.Error()
	case errapi.Authentication:
		errMessage = "invalid credentials"
	default:
		errMessage = "something went wrong"
	}

	return &gqlerror.Error{
		Message:    errMessage,
		Extensions: extensions,
	}
}

func NewExpense(e *expensemodel.Expense) *model.Expense {
	return &model.Expense{
		ID:          e.ID(),
		Description: e.Description(),
		Amount:      e.Amount(),
		Date:        e.Date(),
		UserID:      e.UserID(),
		CreatedAt:   e.CreatedAt(),
		UpdatedAt:   e.UpdatedAt(),
	}
}

func NewPaginatedExpenseResponse(es []*expensemodel.Expense, field string) *model.PaginatedExpenseResponse {

	expenses := make([]*model.Expense, 0)
	for _, e := range es {
		expenses = append(expenses, NewExpense(e))
	}

	cursor := ""
	if len(es) > 0 {
		cursor = utils.BuildCursor(es[len(es)-1], field)
	}

	return &model.PaginatedExpenseResponse{
		Expenses: expenses,
		Cursor:   &cursor,
	}
}
