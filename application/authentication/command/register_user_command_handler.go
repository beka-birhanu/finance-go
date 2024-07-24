package command

import (
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	"github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	timeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	appError "github.com/beka-birhanu/finance-go/application/error"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	"github.com/beka-birhanu/finance-go/domain/model"
)

type UserRegisterCommandHandler struct {
	userRepository repository.IUserRepository
	jwtService     jwt.IJwtService
	hashService    hash.IHashService
	timeService    timeservice.ITimeService
}

// Ensure UserRegisterCommandHandler implements the ICommandHandler interface.
var _ command.ICommandHandler[*UserRegisterCommand, *common.AuthResult] = &UserRegisterCommandHandler{}

func NewRegisterCommandHandler(
	repository repository.IUserRepository,
	jwtService jwt.IJwtService,
	hashService hash.IHashService,
	timeService timeservice.ITimeService,
) *UserRegisterCommandHandler {

	return &UserRegisterCommandHandler{
		userRepository: repository,
		jwtService:     jwtService,
		hashService:    hashService,
		timeService:    timeService,
	}
}

func (h *UserRegisterCommandHandler) Handle(command *UserRegisterCommand) (*common.AuthResult, error) {
	user, err := newUserFromRegisterCommand(command, h.hashService, h.timeService)
	if err != nil {
		return nil, appError.ErrToAppErr(err)
	}

	err = h.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return nil, appError.ErrToAppErr(err)
	}

	return common.NewAuthResult(user.ID(), user.Username(), token), nil
}

func newUserFromRegisterCommand(command *UserRegisterCommand, hashService hash.IHashService, timeService timeservice.ITimeService) (*model.User, error) {
	return model.NewUser(command.Username, command.Password, hashService, timeService.NowUTC())
}