package persistance

import (
	"github.com/beka-birhanu/finance-go/domain/models.go"
)

type IUserRepository interface {
	CreateUser(user *models.User) error
	GetUserById(id string) (*models.User, error)
	GetUserByUsername(id string) (*models.User, error)
	ListUser() ([]*models.User, error)
}
