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

// Handler processes login queries and generates authentication results.
type Handler struct {
	userRepo irepository.IUserRepository
	jwtSvc   ijwt.IService
	hashSvc  hash.IService
}

// Ensure Handler implements iquery.IHandler[*Query, *auth.Result].
var _ iquery.IHandler[*Query, *auth.Result] = &Handler{}

// Config contains dependencies required for creating a new login handler.
type Config struct {
	UserRepository irepository.IUserRepository
	JwtService     ijwt.IService
	HashService    hash.IService
}

// NewHandler initializes a new login handler with the given configuration.
func NewHandler(config Config) *Handler {
	return &Handler{
		userRepo: config.UserRepository,
		jwtSvc:   config.JwtService,
		hashSvc:  config.HashService,
	}
}

// Handle processes a login query and returns an authentication result.
//
// Parameters:
//   - query: A pointer to the Query struct containing username and password.
//
// Returns:
// - *auth.Result: A pointer to the authentication result with user ID, username, and token.
// - error: An error if login fails. Possible errors include:
//   - InvalidCredential: If the username is not found or the password is incorrect.
//   - Unexpected: For unexpected errors during user retrieval or password validation.
func (h *Handler) Handle(query *Query) (*auth.Result, error) {
	user, err := h.userRepo.ByUsername(query.Username)
	if err != nil {
		var domainErr *errdmn.Error
		if errors.As(err, &domainErr) {
			return nil, appError.InvalidCredential(domainErr.Message)
		}
		return nil, errdmn.NewUnexpected(fmt.Sprintf("failed to retrieve user, %v", err))
	}

	isPasswordCorrect, err := h.hashSvc.Match(user.PasswordHash(), query.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to validate user password, %w", err)
	}

	if !isPasswordCorrect {
		return nil, appError.InvalidCredential("incorrect password")
	}

	token, err := h.jwtSvc.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT for user, %w", err)
	}

	return auth.NewResult(user.ID(), user.Username(), token), nil
}

