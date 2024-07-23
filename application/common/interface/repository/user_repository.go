package repository

import (
	"github.com/beka-birhanu/finance-go/domain/model"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserById(id string) (*model.User, error)
	GetUserByUsername(id string) (*model.User, error)
	ListUser() ([]*model.User, error)
}
