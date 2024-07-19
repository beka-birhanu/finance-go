package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/entities"
	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

var users = map[uuid.UUID]entities.User{}

// Ensure UserRepository implements interfaces.persistance.IUserRepository
var _ persistance.IUserRepository = &UserRepository{}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) CreateUser(user *entities.User) error {
	users[user.ID] = *user
	log.Println("userCreated")

	return nil
}

func (u *UserRepository) GetUserById(id string) (*entities.User, error) {
	panic("unimplemented")
}

func (u *UserRepository) GetUserByUsername(username string) (*entities.User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (u *UserRepository) ListUser() ([]*entities.User, error) {
	panic("unimplemented")
}
