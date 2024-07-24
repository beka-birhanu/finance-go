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

var (
	ServerError = New(Unexpected, "server error") // Unexpected server error
)
