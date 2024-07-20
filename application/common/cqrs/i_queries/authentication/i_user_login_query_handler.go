package authentication

import (
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/authentication/queries"
)

type IUserLoginQueryHandler interface {
	Handle(query *queries.UserLoginQuery) (*common.AuthResult, error)
}
