package queries

import (
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
)

type UserLoginQueryHandler struct {
	UserRepository persistance.IUserRepository
	JwtService     jwt.IJwtService
}

func NewUserLoginQueryHandler(repository persistance.IUserRepository, jwtService jwt.IJwtService) *UserLoginQueryHandler {
	return &UserLoginQueryHandler{UserRepository: repository, JwtService: jwtService}
}

func (h *UserLoginQueryHandler) Handle(query *UserLoginQuery) (*common.AuthResult, error) {
	user, err := h.UserRepository.GetUserByUsername(query.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if user.Password != query.Password {
		return nil, fmt.Errorf("invalid username or password")
	}

	token, err := h.JwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return common.NewAuthResult(user.ID, user.Username, token), nil
}
