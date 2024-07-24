/*
Package error provides a way to define and handle application-specific errors.

It provides errors to map specific domain errors to appropriate application error types.
*/
package apperror

import (
	"fmt"
)

// ErrorType represents the type of an error.
type ErrorType string

// Predefined error types.
const (
	// AuthenticationErrorType is used for authentication-related errors.
	Authentication ErrorType = "Authentication"
)

// Error represents a combined application error with a type and message.
type Error struct {
	Type    ErrorType
	Message string
}

// Error formats the Error as a string.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Predefined errors for common scenarios.
var (
	// InvalidCredential represents an authentication error due to invalid credentials.
	InvalidCredential = Error{Type: Authentication, Message: "invalid credentials"}
)

