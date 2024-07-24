package error

import (
	"fmt"

	appError "github.com/beka-birhanu/finance-go/application/error"
)

type StatusCode int

const (
	BadRequestStatusCode     StatusCode = 400 // Bad Request
	ConflictStatusCode       StatusCode = 409 // Conflict
	ServerStatusCode         StatusCode = 500 // Internal Server Error
	AuthenticationStatusCode StatusCode = 401 // Unauthorized
	NotFoundStatusCode       StatusCode = 404 // Not Found
)

type Error struct {
	StatusCode StatusCode
	Message    string
}

func NewErrBadRequest(message string) Error {
	return Error{StatusCode: BadRequestStatusCode, Message: message}
}

func NewErrValidation(message string) Error {
	return Error{StatusCode: BadRequestStatusCode, Message: message}
}

func NewErrConflict(message string) Error {
	return Error{StatusCode: ConflictStatusCode, Message: message}
}

func NewErrServer(message string) Error {
	return Error{StatusCode: ServerStatusCode, Message: message}
}

func NewErrAuthentication(message string) Error {
	return Error{StatusCode: AuthenticationStatusCode, Message: message}
}

func NewErrNotFound(message string) Error {
	return Error{StatusCode: NotFoundStatusCode, Message: message}
}

func NewErrMissingParam(paramName string) Error {
	return Error{StatusCode: BadRequestStatusCode, Message: fmt.Sprintf("%s parameter in path", paramName)}
}

func NewErrInvalidParam(paramName string) Error {
	return Error{StatusCode: BadRequestStatusCode, Message: fmt.Sprintf("invalid %s format", paramName)}
}

func (e Error) Error() string {
	return fmt.Sprint(e.Message)
}

var (
	ErrMissingRequestBody = Error{StatusCode: BadRequestStatusCode, Message: "missing request body"}
)

func ErrToAPIError(err error) Error {
	switch e := err.(type) {
	case appError.Error:
		switch e.Type {
		case appError.ValidationErrorType:
			return Error{StatusCode: BadRequestStatusCode, Message: e.Message}
		case appError.ConflictErrorType:
			return Error{StatusCode: ConflictStatusCode, Message: e.Message}
		case appError.ServerErrorType:
			return Error{StatusCode: ServerStatusCode, Message: e.Message}
		case appError.AuthenticationErrorType:
			return Error{StatusCode: BadRequestStatusCode, Message: "invalid credentials"}
		case appError.NotFoundErrorType:
			return Error{StatusCode: NotFoundStatusCode, Message: e.Message}
		default:
			return Error{StatusCode: ServerStatusCode, Message: "unknown error"}
		}
	default:
		return Error{StatusCode: ServerStatusCode, Message: "something went wrong"}
	}
}

