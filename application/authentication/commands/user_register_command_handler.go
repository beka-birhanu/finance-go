package commands

import (
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/entities"
)

type UserRegisterCommandHandler struct {
	UserRepository persistance.IUserRepository
	JwtService     jwt.IJwtService
}

var (
	ErrUsernameInUse = errors.New("username in use")
)

func NewRegisterCommandHandler(repository persistance.IUserRepository, jwtService jwt.IJwtService) *UserRegisterCommandHandler {
	return &UserRegisterCommandHandler{UserRepository: repository, JwtService: jwtService}
}

func (h *UserRegisterCommandHandler) Handle(command *UserRegisterCommand) (*common.AuthResult, error) {
	if _, err := h.UserRepository.GetUserByUsername(command.Username); err == nil {
		return nil, ErrUsernameInUse
	}

	user := fromRegisterCommand(command)
	if err := h.UserRepository.CreateUser(user); err != nil {
		return nil, fmt.Errorf("server error")
	}

	token, err := h.JwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return common.NewAuthResult(user.ID, user.Username, token), nil
}

func fromRegisterCommand(command *UserRegisterCommand) *entities.User {
	var user entities.User

	user.ID = command.ID
	user.Username = command.Username
	// TODO: hash
	user.Password = command.Password

	return &user
}
