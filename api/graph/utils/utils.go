package utils

import (
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/graph/model"
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
	updatedTime := e.UpdatedAt()
	createdTime := e.CreatedAt()
	return &model.Expense{
		ID:          e.ID(),
		Description: e.Description(),
		Amount:      e.Amount(),
		Date:        e.Date(),
		UserID:      e.UserID(),
		CreatedAt:   &createdTime,
		UpdatedAt:   &updatedTime,
	}
}
