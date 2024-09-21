package utils

import (
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/graph/model"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// NewGQLError creates a new gqlerror with a custom message and extensions
func NewGQLError(err errapi.Error) *gqlerror.Error {
	extensions := make(map[string]interface{})
	extensions["StatusCode"] = err.StatusCode()

	return &gqlerror.Error{
		Message:    err.Error(),
		Extensions: extensions,
	}
}
