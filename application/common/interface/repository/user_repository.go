package irepository

import (
	usermodel "github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/google/uuid"
)

// IUserRepository defines methods for accessing and managing user data.
//
// Methods:
// - Add(user *usermodel.User) error: Creates a new user in the repository and
// returns an error if a conflict occurs.
// - ById(id uuid.UUID) (*usermodel.User, error): Retrieves a user by their unique
// ID and returns the user or an error if the user does not exist.
// - ByUsername(username string) (*usermodel.User, error): Retrieves a user by their
// username and returns the user or an error if the user does not exist.
// - Update(user *usermodel.User) error: Updates the user and returns an error if the user does not exist.
type IUserRepository interface {
	// Add adds a new user to the repository.
	Add(user *usermodel.User) error

	// Update updates the user and returns an error if the user does not exist.
	Update(user *usermodel.User) error

	// ById retrieves a user by their unique identifier.
	ById(id uuid.UUID) (*usermodel.User, error)

	// ByUsername retrieves a user by their username.
	ByUsername(username string) (*usermodel.User, error)
}
