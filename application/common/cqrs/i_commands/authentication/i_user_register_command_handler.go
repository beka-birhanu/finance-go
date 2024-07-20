package authentication

import "github.com/beka-birhanu/finance-go/application/authentication/commands"

type IUserRegisterCommandHandler interface {
	Handle(command *commands.UserRegisterCommand) error
}

