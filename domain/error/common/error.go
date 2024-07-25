package errdmn

import "fmt"

// ErrType represents the type of an error.
type ErrType string

const (
	Validation ErrType = "Validation"  // Validation error
	Conflict   ErrType = "Conflict"    // Conflict error
	Unexpected ErrType = "ServerError" // Unexpected server error
	NotFound   ErrType = "NotFound"    // Resource not found error
)

// Error represents a custom domain error with a type and message.
type Error struct {
	Type    ErrType
	Message string
}

// New creates a new Error with the given type and message.
func New(errType ErrType, message string) *Error {
	return &Error{Type: errType, Message: message}
}

// Error returns the string representation of the Error.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewValidation creates a new validation error with the given message.
func NewValidation(message string) *Error {
	return New(Validation, message)
}

// NewConflict creates a new conflict error with the given message.
func NewConflict(message string) *Error {
	return New(Conflict, message)
}

// NewUnexpected creates a new unexpected server error with the given message.
func NewUnexpected(message string) *Error {
	return New(Unexpected, message)
}

// NewNotFound creates a new not found error with the given message.
func NewNotFound(message string) *Error {
	return New(NotFound, message)
}

