package error

import "fmt"

type ErrorType string

const (
	ValidationErrorType ErrorType = "Validation"
	ConflictErrorType   ErrorType = "Conflict"
	ServerErrorType     ErrorType = "ServerError"
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
	HashingError = Error{Type: ServerErrorType, Message: "error hashing password"}
	ServerError  = Error{Type: ServerErrorType, Message: "unexpected server error"}
)

