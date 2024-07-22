package commands

import (
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/common/authentication"
	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/beka-birhanu/finance-go/domain/entities"
	"github.com/nbutton23/zxcvbn-go"
)

type UserRegisterCommandHandler struct {
	UserRepository persistance.IUserRepository
	JwtService     jwt.IJwtService
	HashService    hash.IHashService
}

var (
	ErrUsernameInUse = domain_errors.UsernameConflict
	ErrWeakPassword  = errors.New("password is too weak!")
)

const (
	MIN_PASSWORD_STRENGTH_SCORE = 3
)

func NewRegisterCommandHandler(repository persistance.IUserRepository, jwtService jwt.IJwtService, hashService hash.IHashService) *UserRegisterCommandHandler {

	return &UserRegisterCommandHandler{UserRepository: repository, JwtService: jwtService, HashService: hashService}
}

func (h *UserRegisterCommandHandler) Handle(command *UserRegisterCommand) (*common.AuthResult, error) {
	user, err := fromRegisterCommand(command, h.HashService)
	if err != nil {
		return nil, err
	}

	err = h.UserRepository.CreateUser(user)
	if errors.Is(err, domain_errors.UsernameConflict) {
		return nil, err
	} else if err != nil {
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

	result := zxcvbn.PasswordStrength(command.Password, nil)
	if result.Score < MIN_PASSWORD_STRENGTH_SCORE {
		return nil, ErrWeakPassword
	}

	hashedPassword, err := hashService.Hash(command.Password)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}
	user.Password = hashedPassword

	return &user, nil
}
