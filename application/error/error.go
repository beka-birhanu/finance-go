package error

import (
	"fmt"

	domainError "github.com/beka-birhanu/finance-go/domain/error"
)

type ErrorType string

const (
	ValidationErrorType     ErrorType = "Validation"
	ConflictErrorType       ErrorType = "Conflict"
	ServerErrorType         ErrorType = "ServerError"
	AuthenticationErrorType ErrorType = "Authentication"
	NotFoundErrorType       ErrorType = "NotFound"
)

type Error struct {
	Type    ErrorType
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

var (
	ErrNegativeExpenseAmount     = Error{Type: ValidationErrorType, Message: "Expense.Amount cannot be negative"}
	ErrExpenseDescriptionTooLong = Error{Type: ValidationErrorType, Message: "Expense.Description is too long"}
	ErrEmptyExpenseDescription   = Error{Type: ValidationErrorType, Message: "Expense.Description cannot be empty"}
)

var (
	ErrUsernameTooShort      = Error{Type: ValidationErrorType, Message: "username is too short"}
	ErrUsernameTooLong       = Error{Type: ValidationErrorType, Message: "username is too long"}
	ErrWeakPassword          = Error{Type: ValidationErrorType, Message: "password is too weak"}
	ErrUsernameInvalidFormat = Error{Type: ValidationErrorType, Message: "username has an invalid format"}
)

var (
	ErrUsernameConflict = Error{Type: ConflictErrorType, Message: "username already taken"}
	ErrIdConflict       = Error{Type: ConflictErrorType, Message: "ID under user and expense don't match"}
)

var (
	ErrInvalidUsernameOrPassword = Error{Type: AuthenticationErrorType, Message: "invalid username or password"}
)

var (
	ErrUserNotFound = Error{Type: NotFoundErrorType, Message: "user not found"}
)
var (
	ServerError = Error{Type: ServerErrorType, Message: "unexpected server error"}
)

func ErrToAppErr(err error) error {
	switch e := err.(type) {
	case domainError.Error:
		switch e.Type {
		case domainError.ValidationErrorType:
			return Error{Type: ValidationErrorType, Message: e.Message}
		case domainError.ConflictErrorType:
			return Error{Type: ConflictErrorType, Message: e.Message}
		case domainError.ServerErrorType:
			return Error{Type: ServerErrorType, Message: e.Message}
		default:
			return Error{Type: ServerErrorType, Message: "unknown error"}
		}

	default:
		return Error{Type: ServerErrorType, Message: "unknown error"}
	}
}
