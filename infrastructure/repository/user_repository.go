package repository

import (
	"database/sql"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	domainError "github.com/beka-birhanu/finance-go/domain/error"
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

var users = map[uuid.UUID]model.User{}
var NotFound = fmt.Errorf("user not found")

// Ensure UserRepository implements interfaces.persistance.IUserRepository
var _ repository.IUserRepository = &UserRepository{}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) CreateUser(user *model.User) error {
	for _, existingUser := range users {
		if existingUser.Username() == user.Username() {
			return domainError.ErrUsernameConflict
		}
	}

	users[user.ID()] = *user

	return nil
}

func (u *UserRepository) GetUserById(id uuid.UUID) (*model.User, error) {

	user, found := users[id]
	if !found {
		return nil, NotFound
	}

	return &user, nil
}

func (u *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	for _, user := range users {
		if user.Username() == username {
			return &user, nil
		}
	}
	return nil, NotFound
}

func (u *UserRepository) ListUser() ([]*model.User, error) {
	var userList []*model.User
	for _, user := range users {
		userList = append(userList, &user)
	}
	return userList, nil
}
