// Package errdmn provides custom error types for the domain layer.
// These error types help in distinguishing between different error scenarios
// such as validation errors, conflicts, unexpected server errors, and not found errors.
package errdmn

import (
	"fmt"

	ierr "github.com/beka-birhanu/finance-go/domain/common/error"
)

// Constants representing different error types.
const (
	Validation = "Validation"  // Validation error type
	Conflict   = "Conflict"    // Conflict error type
	Unexpected = "ServerError" // Unexpected server error type
	NotFound   = "NotFound"    // Resource not found error type
)

// Error represents a custom domain error with a specific type and message.
type Error struct {
	kind    string // The type of the error (e.g., Validation, Conflict)
	Message string // The detailed error message
}

// Ensure that Error implements the ierr.IErr interface.
var _ ierr.IErr = Error{}

// new creates a new Error instance with the given type and message.
func new(errType string, message string) *Error {
	return &Error{kind: errType, Message: message}
}

// Error returns the string representation of the Error.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.kind, e.Message)
}

// Type returns the type of the Error.
func (e Error) Type() string {
	return e.kind
}

// NewValidation creates a new validation error with the given message.
func NewValidation(message string) *Error {
	return new(Validation, message)
}

// NewConflict creates a new conflict error with the given message.
func NewConflict(message string) *Error {
	return new(Conflict, message)
}

// NewUnexpected creates a new unexpected server error with the given message.
func NewUnexpected(message string) *Error {
	return new(Unexpected, message)
}

// NewNotFound creates a new not found error with the given message.
func NewNotFound(message string) *Error {
	return new(NotFound, message)
}
