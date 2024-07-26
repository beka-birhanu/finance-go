package errapi

import (
	"fmt"
)

// StatusCode represents an HTTP status code.
type StatusCode int

// HTTP status codes used in the Error type.
const (
	BadRequest     StatusCode = 400 // Bad Request
	Conflict       StatusCode = 409 // Conflict
	ServerError    StatusCode = 500 // Internal Server Error
	Authentication StatusCode = 401 // Unauthorized
	NotFound       StatusCode = 404 // Not Found
)

// Error represents an API error with a status code and message.
type Error struct {
	statusCode StatusCode
	message    string
}

// NewBadRequest creates a new Error with a 400 Bad Request status code
// and the message provided.
func NewBadRequest(message string) Error {
	return Error{statusCode: BadRequest, message: message}
}

// NewConflict creates a new Error with a 409 Conflict status code
// and the message provided.
func NewConflict(message string) Error {
	return Error{statusCode: Conflict, message: message}
}

// NewServerError creates a new Error with a 500 Internal Server Error status code
// and the message provided.
func NewServerError(message string) Error {
	return Error{statusCode: ServerError, message: message}
}

// NewAuthentication creates a new Error with a 401 Unauthorized status code
// and the message provided.
func NewAuthentication(message string) Error {
	return Error{statusCode: Authentication, message: message}
}

// NewNotFound creates a new Error with a 404 Not Found status code
// and the message provided.
func NewNotFound(message string) Error {
	return Error{statusCode: NotFound, message: message}
}

// Error returns the error message.
func (e Error) Error() string {
	return fmt.Sprint(e.message)
}

