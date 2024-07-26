package registercmd

import (
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	ijwt "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
)

// Handler is a handler for user registration.
type Handler struct {
	userRepository irepository.IUserRepository
	jwtService     ijwt.IService
	hashService    hash.IService
	timeService    itimeservice.IService
}

// Ensure Handler implements the ICommandHandler interface.
var _ icmd.IHandler[*Command, *auth.Result] = &Handler{}

// Config is a configuration struct for creating a new register command handler.
type Config struct {
	UserRepository irepository.IUserRepository
	JwtService     ijwt.IService
	HashService    hash.IService
	TimeService    itimeservice.IService
}

// NewHandler returns a new register command handler using the provided config.
//
// NOTE: All the fields in config are required for the Handle method to function.
func NewHandler(config Config) *Handler {
	return &Handler{
		userRepository: config.UserRepository,
		jwtService:     config.JwtService,
		hashService:    config.HashService,
		timeService:    config.TimeService,
	}
}

// Handle registers a new user from the command received and returns an Result
// if successful.
//
// Returns an error if any of the following happens:
// - The username is already taken by another user.
// - The username does not meet format, length, or validity constraints.
// - The password does not meet the minimum strength requirements.
// - An error occurs during password hashing or generating the JWT.
func (h *Handler) Handle(cmd *Command) (*auth.Result, error) {
	user, err := newUser(cmd, h.hashService, h.timeService)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	err = h.userRepository.Add(user)
	if err != nil {
		return nil, fmt.Errorf("failed to add new user to repository: %w", err)
	}

	token, err := h.jwtService.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT for user: %w", err)
	}

	return auth.NewResult(user.ID(), user.Username(), token), nil
}

// newUser creates a new user instance using the provided command, hash service,
// and time service.
//
// Returns an error if user creation fails.
func newUser(cmd *Command, hashService hash.IService, timeService itimeservice.IService) (*usermodel.User, error) {
	config := usermodel.Config{
		Username:       cmd.Username,
		PlainPassword:  cmd.Password,
		CreationTime:   timeService.NowUTC(),
		PasswordHasher: hashService,
	}
	return usermodel.New(config)
}
