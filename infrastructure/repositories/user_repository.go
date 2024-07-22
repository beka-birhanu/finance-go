package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/beka-birhanu/finance-go/domain/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	DB *sql.DB
}

var users = map[uuid.UUID]models.User{}
var NotFound = errors.New("username already taken")

// Ensure UserRepository implements interfaces.persistance.IUserRepository
var _ persistance.IUserRepository = &UserRepository{}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (u *UserRepository) CreateUser(user *models.User) error {
	for _, existingUser := range users {
		if existingUser.Username == user.Username {
			return domain_errors.ErrUsernameConflict
		}
	}

	users[user.ID] = *user

	return nil
}

func (u *UserRepository) GetUserById(id string) (*models.User, error) {
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

func (u *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, NotFound
}

func (u *UserRepository) ListUser() ([]*models.User, error) {
	var userList []*models.User
	for _, user := range users {
		userList = append(userList, &user)
	}
	return userList, nil
}
