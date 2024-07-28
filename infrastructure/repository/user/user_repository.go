package userrepo

import (
	"database/sql"

	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	erruser "github.com/beka-birhanu/finance-go/domain/error/user"
	"github.com/beka-birhanu/finance-go/domain/model/user"
	expenserepo "github.com/beka-birhanu/finance-go/infrastructure/repository/expense"
	"github.com/google/uuid"
)

// Repository handles the persistence of user models.
type Repository struct {
	db *sql.DB
}

var users = map[uuid.UUID]usermodel.User{} // In-memory storage for users

// Ensure UserRepository implements repository.IUserRepository.
var _ irepository.IUserRepository = &Repository{}

// New creates a new UserRepository with the given database connection.
func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Save inserts or updates a user in the repository.
// If the user already exists, it updates the existing record.
// If the user does not exist, it adds a new record.
//
// Returns:
//   - error: An error if a conflict occurs, otherwise nil.
func (u *Repository) Save(user *usermodel.User) error {
	_, found := users[user.ID()]
	if !found {
		for _, existingUser := range users {
			if existingUser.Username() == user.Username() {
				return erruser.UsernameConflict
			}
		}
	}

	users[user.ID()] = *user
	for _, expense := range user.Expenses() {
		key := expenserepo.PrimaryKey{
			Id:     expense.ID(),
			UserId: user.ID(),
		}
		expenserepo.Expenses[key] = expense
	}
	return nil
}

// ById retrieves a user by their ID.
//
// Returns:
//   - *usermodel.User: A pointer to the retrieved user model.
//   - error: An error if the user is not found, otherwise nil.
func (u *Repository) ById(id uuid.UUID) (*usermodel.User, error) {
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
func (u *Repository) ByUsername(username string) (*usermodel.User, error) {
	for _, user := range users {
		if user.Username() == username {
			return &user, nil
		}
	}
	return nil, erruser.NotFound
}
