package query

import (
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	"github.com/beka-birhanu/finance-go/application/common/interface/jwt"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	appError "github.com/beka-birhanu/finance-go/application/error"
	"github.com/beka-birhanu/finance-go/domain/common/hash"
)

type UserLoginQueryHandler struct {
	userRepository repository.IUserRepository
	jwtService     jwt.IJwtService
	hashService    hash.IHashService
}

var _ query.IQueryHandler[*UserLoginQuery, *common.AuthResult] = &UserLoginQueryHandler{}

func NewUserLoginQueryHandler(repository repository.IUserRepository, jwtService jwt.IJwtService, hashService hash.IHashService) *UserLoginQueryHandler {
	return &UserLoginQueryHandler{userRepository: repository, jwtService: jwtService, hashService: hashService}
}

func (h *UserLoginQueryHandler) Handle(query *UserLoginQuery) (*common.AuthResult, error) {
	user, err := h.userRepository.GetUserByUsername(query.Username)
	if err != nil {
		return nil, appError.ErrInvalidUsernameOrPassword
	}

	isPassowrdCorrect, err := h.hashService.Match(user.PasswordHash(), query.Password)
	if err != nil {
		return nil, appError.ServerError
	}

	if !isPassowrdCorrect {
		return nil, appError.ErrInvalidUsernameOrPassword
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return nil, appError.ServerError
	}

	return common.NewAuthResult(user.ID(), user.Username(), token), nil
}
