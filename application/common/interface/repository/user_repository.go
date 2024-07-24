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
// - CreateUser(user *usermodel.User) error: Creates a new user in the repository and
// returns an error if any occurs.
// - GetUserById(id uuid.UUID) (*usermodel.User, error): Retrieves a user by their unique
// ID and returns the user or an error if the user does not exist.
// - GetUserByUsername(username string) (*usermodel.User, error): Retrieves a user by their
// username and returns the user or an error if the user does not exist.
// - ListUser() ([]*usermodel.User, error): Lists all users in the repository and returns
// a slice of users or an error.
type IUserRepository interface {
	// CreateUser adds a new user to the repository.
	CreateUser(user *usermodel.User) error

	// GetUserById retrieves a user by their unique identifier.
	GetUserById(id uuid.UUID) (*usermodel.User, error)

	// GetUserByUsername retrieves a user by their username.
	GetUserByUsername(username string) (*usermodel.User, error)

	// ListUser returns a list of all users in the repository.
	ListUser() ([]*usermodel.User, error)
}

