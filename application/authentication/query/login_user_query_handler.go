// Package loginqry provides functionality for handling login queries.
package loginqry

import (
	"errors"
	"fmt"

	auth "github.com/beka-birhanu/finance-go/application/authentication/common"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	ijwt "github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	appError "github.com/beka-birhanu/finance-go/application/error"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
)

// Handler processes login queries.
type Handler struct {
	userRepository irepository.IUserRepository
	jwtService     ijwt.IJwtService
	hashService    hash.IHashService
}

// Ensuring that Handler implements iquery.IHandler[*Query, *auth.Result]
var _ iquery.IHandler[*Query, *auth.Result] = &Handler{}

// Config holds the configuration for creating a new Handler.
type Config struct {
	UserRepository irepository.IUserRepository
	JwtService     ijwt.IJwtService
	HashService    hash.IHashService
}

// NewHandler creates a new Handler with the provided configuration.
func NewHandler(config Config) *Handler {
	return &Handler{
		userRepository: config.UserRepository,
		jwtService:     config.JwtService,
		hashService:    config.HashService,
	}
}

// Handle processes the login query and returns an auth.Result.
//
// Parameters:
//   - query: A pointer to the Query containing the username and password.
//
// Returns:
// - *auth.Result: A pointer to the auth.Result containing the user ID, username, and token.
// - error: An error if the login process fails. Possible errors include:
//   - InvalidCredential: If the username does not exist or the password is incorrect.
//   - Unexpected: If there is an unexpected error during user retrieval or password validation.
func (h *Handler) Handle(query *Query) (*auth.Result, error) {
	user, err := h.userRepository.ByUsername(query.Username)
	if err != nil {
		var domainErr *errdmn.Error
		if errors.As(err, &domainErr) {
			return nil, appError.InvalidCredential(domainErr.Message)
		}
		return nil, errdmn.NewUnexpected(fmt.Sprintf("failed to retrieve user, %v", err))
	}

	isPasswordCorrect, err := h.hashService.Match(user.PasswordHash(), query.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to validate user password, %w", err)
	}

	if !isPasswordCorrect {
		return nil, appError.InvalidCredential("incorrect password")
	}

	token, err := h.jwtService.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT for user, %w", err)
	}

	return auth.NewResult(user.ID(), user.Username(), token), nil
}
