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

// Handler is responsible for user registration.
type Handler struct {
	userRepo irepository.IUserRepository
	jwtSvc   ijwt.IService
	hashSvc  hash.IService
	timeSvc  itimeservice.IService
}

// Ensure Handler implements the ICommandHandler interface.
var _ icmd.IHandler[*Command, *auth.Result] = &Handler{}

// Config holds the dependencies needed to create a new Handler.
type Config struct {
	UserRepo irepository.IUserRepository
	JwtSvc   ijwt.IService
	HashSvc  hash.IService
	TimeSvc  itimeservice.IService
}

// NewHandler creates a new Handler with the provided configuration.
func NewHandler(cfg Config) *Handler {
	return &Handler{
		userRepo: cfg.UserRepo,
		jwtSvc:   cfg.JwtSvc,
		hashSvc:  cfg.HashSvc,
		timeSvc:  cfg.TimeSvc,
	}
}

// Handle processes a user registration command and returns a result if successful.
// It returns an error for issues such as:
// - Username already taken
// - Invalid username format
// - Weak password
// - Errors during hashing or JWT generation
func (h *Handler) Handle(cmd *Command) (*auth.Result, error) {
	user, err := createUser(cmd, h.hashSvc, h.timeSvc)
	if err != nil {
		return nil, fmt.Errorf("creating new user failed: %w", err)
	}

	if err := h.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("saving user to repository failed: %w", err)
	}

	token, err := h.jwtSvc.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("JWT generation failed: %w", err)
	}

	return auth.NewResult(user.ID(), user.Username(), token), nil
}

// createUser initializes a new user instance using the provided command, hash service,
// and time service. It returns an error if user creation fails.
func createUser(cmd *Command, hashSvc hash.IService, timeSvc itimeservice.IService) (*usermodel.User, error) {
	cfg := usermodel.Config{
		Username:       cmd.Username,
		PlainPassword:  cmd.Password,
		CreationTime:   timeSvc.NowUTC(),
		PasswordHasher: hashSvc,
	}
	return usermodel.New(cfg)
}
