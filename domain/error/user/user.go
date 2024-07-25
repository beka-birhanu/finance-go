/*
Package erruser defines user-related errors for the application.

It provides a set of predefined errors related to user not-found, validation,
conflicts, and unexpected issues. These errors are used throughout the application
to handle various error conditions specific to user operations.
*/
package erruser

import "github.com/beka-birhanu/finance-go/domain/error/common"

// Validation errors
var (
	// Username is shorter than allowed.
	UsernameTooShort = errdmn.NewValidation("username is too short.")

	// Username is longer than allowed.
	UsernameTooLong = errdmn.NewValidation("username is too long.")

	// Password is susceptible to attack.
	WeakPassword = errdmn.NewValidation("password is too weak.")

	// Username is not UUID.
	UsernameInvalidFormat = errdmn.NewValidation("username has an invalid format.")
)

// Conflict errors
var (
	// User with a similar username exists.
	UsernameConflict = errdmn.NewConflict("username already taken.")

	// ID under user and expense do not match.
	ExpenseIdConflict = errdmn.NewConflict("ID under user and expense do not match.")
)

// NotFound errors
var (
	// User is does not exist.
	NotFound = errdmn.NewNotFound("User not found.")
)

// Unexpected errors
var (
	// Unexpected error occurred while hashing.
	Hash = errdmn.NewUnexpected("error hashing password.")
)
