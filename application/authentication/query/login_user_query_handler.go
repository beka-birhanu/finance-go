package query

import (
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
)

type UserLoginQueryHandler struct {
	userRepository repository.IUserRepository
	jwtService     jwt.IJwtService
	hashService    hash.IHashService
}

func NewUserLoginQueryHandler(repository repository.IUserRepository, jwtService jwt.IJwtService, hashService hash.IHashService) *UserLoginQueryHandler {
	return &UserLoginQueryHandler{userRepository: repository, jwtService: jwtService, hashService: hashService}
}

func (h *UserLoginQueryHandler) Handle(query *UserLoginQuery) (*common.AuthResult, error) {
	user, err := h.userRepository.GetUserByUsername(query.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	isPassowrdCorrect, err := h.hashService.Match(user.PasswordHash(), query.Password)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	if !isPassowrdCorrect {
		return nil, fmt.Errorf("invalid username or password")
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return common.NewAuthResult(user.ID(), user.Username(), token), nil
}
