package persistance

import "github.com/beka-birhanu/finance-go/domain/entities"

type IUserRepository interface {
	CreateUser(user *entities.User) error
	GetUserById(id string) (*entities.User, error)
	GetUserByUsername(id string) (*entities.User, error)
	ListUser() ([]*entities.User, error)
}
