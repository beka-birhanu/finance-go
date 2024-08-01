/*
Package usermodel defines the `User` aggregate, which represents an individual user,
and includes methods for creating and managing users. It handles user creation,
validation of usernames and passwords, and association with expenses.

Key Components:
  - User: Represents a user with details such as username, password hash, and associated
    expenses.
  - Config: Holds the mandatory parameters required to create a new User.
  - New: Creates a new User instance based on the provided configuration.

Dependencies:
- github.com/google/uuid: Used for generating unique IDs.
- github.com/nbutton23/zxcvbn-go: Used for password strength evaluation.
- time: Used for timestamps.
*/
package usermodel

import (
	"regexp"
	"time"

	"github.com/beka-birhanu/finance-go/domain/common/hash"
	"github.com/beka-birhanu/finance-go/domain/error/user"
	"github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
	"github.com/nbutton23/zxcvbn-go"
)

const (
	minPasswordStrengthScore = 3

	usernamePattern   = `^[a-zA-Z0-9_]+$` // Alphanumeric with underscores
	minUsernameLength = 3
	maxUsernameLength = 20
)

var (
	usernameRegex = regexp.MustCompile(usernamePattern)
)

// User represents the Aggregate user.
type User struct {
	id           uuid.UUID
	username     string
	passwordHash string
	createdAt    time.Time
	updatedAt    time.Time
	expenses     []expensemodel.Expense
}

// Config holds all mandatory parameters for creating a new User.
type Config struct {
	// Username must be non-empty and adhere to the username format and length constraints.
	Username string

	// PlainPassword must meet the minimum password strength requirements.
	PlainPassword string

	// CreationTime is the timestamp when the User is created.
	CreationTime time.Time

	// PasswordHasher is a service used to hash the plain password.
	PasswordHasher hash.IService
}

// ConfigForExistingHash holds all parameters for creating a User with an existing password hash.
type ConfigForExistingHash struct {
	ID           uuid.UUID // Unique identifier for the user
	Username     string    // Username of the user
	PasswordHash string    // Pre-hashed password for the user
	CreationTime time.Time // Timestamp when the user was created
	UpdatedAt    time.Time // Timestamp when the user was last updated
}

// New creates a new User with the provided configuration.
//
// Returns:
// - A pointer to the newly created User if successful.
// - An error if any of the following conditions are not met:
//   - Any field in the config is missing.
//   - The username does not meet format, length, or validity constraints.
//   - The password does not meet the minimum strength requirements.
//   - An error occurs during password hashing.
func New(config Config) (*User, error) {
	if err := validateUsername(config.Username); err != nil {
		return nil, err
	}

	if err := validatePassword(config.PlainPassword); err != nil {
		return nil, err
	}

	passwordHash, err := config.PasswordHasher.Hash(config.PlainPassword)
	if err != nil {
		return nil, erruser.Hash
	}

	return &User{
		id:           uuid.New(), // New ID for the user
		username:     config.Username,
		passwordHash: passwordHash,
		createdAt:    config.CreationTime,
		updatedAt:    config.CreationTime,
		expenses:     []expensemodel.Expense{}, // Ensure slice is initialized
	}, nil
}

// NewWithExistingHash creates a new User with the provided configuration, where the password is already hashed.
//
// Returns:
// - A pointer to the newly created User if successful.
// - An error if any of the following conditions are not met:
//   - The username does not meet format, length, or validity constraints.
//   - The password hash is not valid or empty.
//   - Any other unexpected error occurs during user creation.
func NewWithExistingHash(config ConfigForExistingHash) (*User, error) {
	if err := validateUsername(config.Username); err != nil {
		return nil, err
	}

	return &User{
		id:           config.ID,
		username:     config.Username,
		passwordHash: config.PasswordHash,
		createdAt:    config.CreationTime,
		updatedAt:    config.UpdatedAt,
		expenses:     []expensemodel.Expense{}, // Ensure slice is initialized
	}, nil
}

// validateUsername validates the username according to the defined rules.
func validateUsername(username string) error {
	if len(username) < minUsernameLength {
		return erruser.UsernameTooShort
	}
	if len(username) > maxUsernameLength {
		return erruser.UsernameTooLong
	}
	if !usernameRegex.MatchString(username) {
		return erruser.UsernameInvalidFormat
	}
	return nil
}

// validatePassword checks the strength of the password.
func validatePassword(password string) error {
	result := zxcvbn.PasswordStrength(password, nil)
	if result.Score < minPasswordStrengthScore {
		return erruser.WeakPassword
	}
	return nil
}

// ID returns the user's ID.
func (u *User) ID() uuid.UUID {
	return u.id
}

// Username returns the user's username.
func (u *User) Username() string {
	return u.username
}

// PasswordHash returns the user's password hash.
func (u *User) PasswordHash() string {
	return u.passwordHash
}

// CreatedAt returns the user's creation timestamp.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the user's last updated timestamp.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// Expenses returns a copy of the user's expenses.
func (u *User) Expenses() []expensemodel.Expense {
	expensesCopy := make([]expensemodel.Expense, len(u.expenses))
	copy(expensesCopy, u.expenses)
	return expensesCopy
}

// AddExpense adds an expense to the user and updates the user's last updated timestamp.
// It ensures that the expense's user ID matches the user's ID.
//
// Parameters:
// - currentUTCTime: The current UTC time, used to update the user's last updated timestamp.
//
// Returns:
// - An error if the expense's UserID does not match the user's ID, otherwise returns nil.
func (u *User) AddExpense(expense *expensemodel.Expense, currentUTCTime time.Time) error {
	if expense.UserID() != u.id {
		return erruser.ExpenseIdConflict
	}

	copyExpense := *expense
	u.expenses = append(u.expenses, copyExpense)
	u.updatedAt = currentUTCTime
	return nil
}

