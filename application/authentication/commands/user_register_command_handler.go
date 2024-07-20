package commands

import (
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/hash"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/entities"
)

type UserRegisterCommandHandler struct {
	UserRepository persistance.IUserRepository
	JwtService     jwt.IJwtService
	HashService    hash.IHashService
}

var (
	ErrUsernameInUse = errors.New("username in use")
)

func NewRegisterCommandHandler(repository persistance.IUserRepository, jwtService jwt.IJwtService, hashService hash.IHashService) *UserRegisterCommandHandler {

	return &UserRegisterCommandHandler{UserRepository: repository, JwtService: jwtService, HashService: hashService}
}

func (h *UserRegisterCommandHandler) Handle(command *UserRegisterCommand) (*common.AuthResult, error) {
	if _, err := h.UserRepository.GetUserByUsername(command.Username); err == nil {
		return nil, ErrUsernameInUse
	}

	user, err := fromRegisterCommand(command, h.HashService)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	if err := h.UserRepository.CreateUser(user); err != nil {
		return nil, fmt.Errorf("server error")
	}

	token, err := h.JwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return common.NewAuthResult(user.ID, user.Username, token), nil
}

func fromRegisterCommand(command *UserRegisterCommand, hashService hash.IHashService) (*entities.User, error) {
	var user entities.User

	user.ID = command.ID
	user.Username = command.Username

	hashedPassword, err := hashService.Hash(command.Password)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}
	user.Password = hashedPassword

	return &user, nil
}
