package command

import (
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	"github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	domainError "github.com/beka-birhanu/finance-go/domain/error"
	"github.com/beka-birhanu/finance-go/domain/model"
)

type UserRegisterCommandHandler struct {
	userRepository repository.IUserRepository
	jwtService     jwt.IJwtService
	hashService    hash.IHashService
}

// Ensure UserRegisterCommandHandler implements the ICommandHandler interface.
var _ command.ICommandHandler[*UserRegisterCommand, *common.AuthResult] = &UserRegisterCommandHandler{}

func NewRegisterCommandHandler(repository repository.IUserRepository, jwtService jwt.IJwtService, hashService hash.IHashService) *UserRegisterCommandHandler {
	return &UserRegisterCommandHandler{userRepository: repository, jwtService: jwtService, hashService: hashService}
}

func (h *UserRegisterCommandHandler) Handle(command *UserRegisterCommand) (*common.AuthResult, error) {
	user, err := fromRegisterCommand(command, h.hashService)
	if err != nil {
		return nil, err
	}

	err = h.userRepository.CreateUser(user)
	if errors.Is(err, domainError.ErrUsernameConflict) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("server error")
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return common.NewAuthResult(user.ID(), user.Username(), token), nil
}

func fromRegisterCommand(command *UserRegisterCommand, hashService hash.IHashService) (*model.User, error) {
	return model.NewUser(command.Username, command.Password, hashService)
}
