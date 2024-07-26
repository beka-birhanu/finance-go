package repository

import (
	"database/sql"

	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	"github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/google/uuid"
)

// UserRepository handles the persistence of user models.
type UserRepository struct {
	DB *sql.DB
}

var users = map[uuid.UUID]usermodel.User{} // In-memory storage for users

// Ensure UserRepository implements repository.IUserRepository.
var _ irepository.IUserRepository = &UserRepository{}

// NewUserRepository creates a new UserRepository with the given database connection.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// Add inserts a new user into the repository.
// Returns:
//   - error: An error if the user already exists, otherwise nil.
func (u *UserRepository) Add(user *usermodel.User) error {
	for _, existingUser := range users {
		if existingUser.Username() == user.Username() {
			return erruser.UsernameConflict
		}
	}

	users[user.ID()] = *user
	return nil
}

// ById retrieves a user by their ID.
//
// Returns:
//   - *usermodel.User: A pointer to the retrieved user model.
//   - error: An error if the user is not found, otherwise nil.
func (u *UserRepository) ById(id uuid.UUID) (*usermodel.User, error) {
	user, found := users[id]
	if !found {
		return nil, erruser.NotFound
	}

	return &user, nil
}

// ByUsername retrieves a user by their username.
//
// Returns:
//   - *usermodel.User: A pointer to the retrieved user model.
//   - error: An error if the user is not found, otherwise nil.
func (u *UserRepository) ByUsername(username string) (*usermodel.User, error) {
	for _, user := range users {
		if user.Username() == username {
			return &user, nil
		}
	}
	return nil, erruser.NotFound
}

// Update modifies an existing user in the repository.
//
// Returns:
//   - error: An error if the user is not found, otherwise nil.
func (u *UserRepository) Update(user *usermodel.User) error {
	_, err := u.ById(user.ID())
	if err != nil {
		return err
	}

	users[user.ID()] = *user
	return nil
}

