/*
Package repository provides an interface for user data access operations.

It includes the `IUserRepository` interface for managing user data in a repository.
*/
package irepository

import (
	"github.com/beka-birhanu/finance-go/domain/model/user"
	"github.com/google/uuid"
)

// IUserRepository defines methods for accessing and managing user data.
//
// Methods:
// - AddUser(user *usermodel.User) error: Creates a new user in the repository and
// returns an error if conflict occurs.
// - ById(id uuid.UUID) (*usermodel.User, error): Retrieves a user by their unique
// ID and returns the user or an error if the user does not exist.
// - ByUsername(username string) (*usermodel.User, error): Retrieves a user by their
// username and returns the user or an error if the user does not exist.
// - Update(user *usermodel.User) error: updates the user passed and returns an error if user doesnot exist
type IUserRepository interface {
	// Add adds a new user to the repository.
	Add(user *usermodel.User) error

	// Update updates the user passed and returns an error if user doesnot exist
	Update(user *usermodel.User) error

	// UserById retrieves a user by their unique identifier.
	ById(id uuid.UUID) (*usermodel.User, error)

	// GetUserByUsername retrieves a user by their username.
	ByUsername(username string) (*usermodel.User, error)
}
