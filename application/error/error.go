/*
Package error provides a way to define and handle application-specific errors.

It provides errors to map specific domain errors to appropriate application error types.
*/
package apperror

import (
	"fmt"

	ierr "github.com/beka-birhanu/finance-go/domain/common/error"
)

// Predefined error types.
const (
	// AuthenticationErrorType is used for authentication-related errors.
	Authentication = "Authentication"
)

// Error represents a combined application error with a type and message.
type Error struct {
	kind    string
	Message string
}

var _ ierr.IErr = Error{} // Making sure Error implements IErr

// Error formats the Error as a string.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.kind, e.Message)
}

// Type returns the type of the Error.
func (e Error) Type() string {
	return e.kind
}

// Predefined errors for common scenarios.
var (
	// InvalidCredential represents an authentication error due to invalid credentials.
	InvalidCredential = Error{kind: Authentication, Message: "invalid credentials"}
)
