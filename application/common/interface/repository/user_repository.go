package repository

import (
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/google/uuid"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserById(id uuid.UUID) (*model.User, error)
	GetUserByUsername(id string) (*model.User, error)
	ListUser() ([]*model.User, error)
}
