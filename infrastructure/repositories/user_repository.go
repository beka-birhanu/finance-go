package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/beka-birhanu/finance-go/domain/entities"
	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

var users = map[uuid.UUID]entities.User{}
var NotFound = errors.New("username already taken")

// Ensure UserRepository implements interfaces.persistance.IUserRepository
var _ persistance.IUserRepository = &UserRepository{}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) CreateUser(user *entities.User) error {
	for _, existingUser := range users {
		if existingUser.Username == user.Username {
			return domain_errors.UsernameConflict
		}
	}

	users[user.ID] = *user

	return nil
}

func (u *UserRepository) GetUserById(id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format: %v", err)
	}

	user, found := users[userID]
	if !found {
		return nil, NotFound
	}

	return &user, nil
}

func (u *UserRepository) GetUserByUsername(username string) (*entities.User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, NotFound
}

func (u *UserRepository) ListUser() ([]*entities.User, error) {
	var userList []*entities.User
	for _, user := range users {
		userList = append(userList, &user)
	}
	return userList, nil
}

