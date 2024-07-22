package queries

import (
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/common/authentication"
)

type UserLoginQueryHandler struct {
	UserRepository persistance.IUserRepository
	JwtService     jwt.IJwtService
	HashService    hash.IHashService
}

func NewUserLoginQueryHandler(repository persistance.IUserRepository, jwtService jwt.IJwtService, hashService hash.IHashService) *UserLoginQueryHandler {
	return &UserLoginQueryHandler{UserRepository: repository, JwtService: jwtService, HashService: hashService}
}

func (h *UserLoginQueryHandler) Handle(query *UserLoginQuery) (*common.AuthResult, error) {
	user, err := h.UserRepository.GetUserByUsername(query.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	isPassowrdCorrect, err := h.HashService.Match(user.PasswordHash, query.Password)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	if !isPassowrdCorrect {
		return nil, fmt.Errorf("invalid username or password")
	}

	token, err := h.JwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return common.NewAuthResult(user.ID, user.Username, token), nil
}
