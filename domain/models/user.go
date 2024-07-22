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
	ID           uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}, nil
}

func (u *User) GetExpenses() *[]Expense {
	expensesCopy := make([]Expense, len(u.expenses))
	copy(expensesCopy, u.expenses)
	return &expensesCopy
}

func (u *User) AddExpense(expense *Expense) error {
	if expense.ID() != u.ID {
		return fmt.Errorf("ID under user and expense dont match")
	}
	copyExpense := *expense
	u.expenses = append(u.expenses, copyExpense)
	return nil
}
