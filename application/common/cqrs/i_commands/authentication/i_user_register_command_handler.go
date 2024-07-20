package authentication

import (
	"github.com/beka-birhanu/finance-go/application/authentication/commands"
	"github.com/beka-birhanu/finance-go/application/authentication/common"
)

type IUserRegisterCommandHandler interface {
	Handle(command *commands.UserRegisterCommand) (*common.AuthResult, error)
}
