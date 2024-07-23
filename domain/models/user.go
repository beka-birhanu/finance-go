package models

import (
	"fmt"
	"regexp"
	"time"

	hash "github.com/beka-birhanu/finance-go/domain/common/authentication"
	"github.com/beka-birhanu/finance-go/domain/domain_errors"
	"github.com/google/uuid"
	"github.com/nbutton23/zxcvbn-go"
)

const (
	MIN_PASSWORD_STRENGTH_SCORE = 3

	USERNAME_PATTERN    = `^[a-zA-Z0-9_]+$` // Alphanumeric with underscores
	MIN_USERNAME_LENGTH = 3
	MAX_USERNAME_LENGTH = 20
)

var (
	usernameRegex = regexp.MustCompile(USERNAME_PATTERN)
)

type User struct {
	id           uuid.UUID
	username     string
	passwordHash string
	createdAt    time.Time
	updatedAt    time.Time
	expenses     []Expense
}

func NewUser(username, plainPassword string, passwordHasher hash.IHashService) (*User, error) {
	if len(username) < MIN_USERNAME_LENGTH {
		return nil, domain_errors.ErrUsernameTooShort
	}
	if len(username) > MAX_USERNAME_LENGTH {
		return nil, domain_errors.ErrUsernameTooLong
	}

	if !usernameRegex.MatchString(username) {
		return nil, domain_errors.ErrUsernameInvalidFormat
	}

	result := zxcvbn.PasswordStrength(plainPassword, nil)
	if result.Score < MIN_PASSWORD_STRENGTH_SCORE {
		return nil, domain_errors.ErrWeakPassword
	}

	passwordHash, err := passwordHasher.Hash(plainPassword)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	return &User{
		id:           uuid.New(),
		username:     username,
		passwordHash: passwordHash,
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
		expenses:     []Expense{},
	}, nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) Expenses() []Expense {
	expensesCopy := make([]Expense, len(u.expenses))
	copy(expensesCopy, u.expenses)
	return expensesCopy
}

func (u *User) AddExpense(expense *Expense) error {
	if expense.UserID() != u.id {
		return fmt.Errorf("ID under user and expense don't match")
	}
	copyExpense := *expense
	u.expenses = append(u.expenses, copyExpense)
	return nil
}
