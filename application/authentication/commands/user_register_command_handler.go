package commands

import (
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/entities"
)

type UserRegisterCommandHandler struct {
	UserRepository persistance.IUserRepository
}

var (
	ErrUsernameInUse = errors.New("username in use")
)

func NewRegisterCommandHandler(repository persistance.IUserRepository) *UserRegisterCommandHandler {
	return &UserRegisterCommandHandler{UserRepository: repository}
}

func (h *UserRegisterCommandHandler) Handle(command *UserRegisterCommand) error {
	if _, err := h.UserRepository.GetUserByUsername(command.Username); err == nil {
		return ErrUsernameInUse
	}

	user := fromRegisterCommand(command)
	if err := h.UserRepository.CreateUser(user); err != nil {
		return fmt.Errorf("server error")
	}

	return nil
}

func fromRegisterCommand(command *UserRegisterCommand) *entities.User {
	var user entities.User

	user.ID = command.ID
	user.Username = command.Username
	// TODO: hash
	user.Password = command.Password

	return &user
}
