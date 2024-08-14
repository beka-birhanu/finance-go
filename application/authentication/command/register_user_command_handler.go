// Package registercmd provides the command handler for user registration.
// It processes user registration commands and handles user creation,
// password hashing, and JWT token generation.
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

// Handler is responsible for processing user registration commands.
// It interacts with the user repository, JWT service, hash service,
// and time service to complete the registration process.
type Handler struct {
	userRepo irepository.IUserRepository
	jwtSvc   ijwt.IService
	hashSvc  hash.IService
	timeSvc  itimeservice.IService
}

// Ensure Handler implements the ICommandHandler interface for Command type and Result type.
var _ icmd.IHandler[*Command, *auth.Result] = &Handler{}

// Config holds the dependencies needed to create a new Handler.
type Config struct {
	UserRepo irepository.IUserRepository
	JwtSvc   ijwt.IService
	HashSvc  hash.IService
	TimeSvc  itimeservice.IService
}

// NewHandler creates a new Handler with the provided configuration.
// It initializes the Handler with the necessary services for user registration.
func NewHandler(cfg Config) *Handler {
	return &Handler{
		userRepo: cfg.UserRepo,
		jwtSvc:   cfg.JwtSvc,
		hashSvc:  cfg.HashSvc,
		timeSvc:  cfg.TimeSvc,
	}
}

// Handle processes a user registration command and returns an authentication result if successful.
// Returns an error if any of the following occur:
// - Username is already taken.
// - Invalid username format.
// - Weak password.
// - Errors during user creation, saving, or JWT generation.
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

// createUser initializes a new user instance using the provided command,
// hash service, and time service. It returns an error if user creation fails.
// This function creates a user with hashed password and current creation time.
func createUser(cmd *Command, hashSvc hash.IService, timeSvc itimeservice.IService) (*usermodel.User, error) {
	cfg := usermodel.Config{
		Username:       cmd.Username,
		PlainPassword:  cmd.Password,
		CreationTime:   timeSvc.NowUTC(),
		PasswordHasher: hashSvc,
	}
	return usermodel.New(cfg)
}
