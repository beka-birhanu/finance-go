package queries

import (
	"fmt"

	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
)

type UserLoginQueryHandler struct {
	UserRepository persistance.IUserRepository
}

func NewUserLogingQueryHandler(repository persistance.IUserRepository) *UserLoginQueryHandler {
	return &UserLoginQueryHandler{UserRepository: repository}
}

func (h *UserLoginQueryHandler) Handle(query *UserLoginQuery) (*common.AuthResult, error) {
	user, err := h.UserRepository.GetUserByUsername(query.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	if user.Password != query.Password {
		return nil, fmt.Errorf("invalid username or password")
	}

	return common.NewAuthResult(user.ID, query.Username, "token"), nil
}
