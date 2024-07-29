package baseapi

import (
	"errors"
	"fmt"
	"net/http"

	errapi "github.com/beka-birhanu/finance-go/api/error"
	httputil "github.com/beka-birhanu/finance-go/api/http_util"
	"github.com/go-playground/validator/v10"
)

type BaseHandler struct{}

// ValidatedBody validates the body of an HTTP request and populate the given interface.
func (h *BaseHandler) ValidatedBody(r *http.Request, s interface{}) error {
	// Parse JSON body
	if err := httputil.ParseJSON(r, s); err != nil {
		var apiErr errapi.Error
		if errors.As(err, &apiErr) {
			return apiErr
		}
		return errapi.NewServerError("unable to parse request body")
	}

	// Validate the parsed struct
	if err := httputil.Validate.Struct(s); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return errapi.NewBadRequest(fmt.Sprintf("invalid payload: %v", validationErrs))
		}
		return errapi.NewServerError("request body validation failed")
	}

	return nil
}

func (h *BaseHandler) Problem(w http.ResponseWriter, err errapi.Error) {
	// shadow error messages if they are not supposed to be seen.
	switch err.StatusCode() {
	case errapi.BadRequest:
		httputil.RespondError(w, err)

	case errapi.Conflict:
		httputil.RespondError(w, err)

	case errapi.NotFound:
		httputil.RespondError(w, err)

	// client should not know what coused the fail
	case errapi.Authentication:
		shadowedErr := errapi.NewAuthentication("invalid credentials")
		httputil.RespondError(w, shadowedErr)

	// client should not know what coused the fail
	case errapi.ServerError:
		shadowedErr := errapi.NewServerError("something went wrong")
		httputil.RespondError(w, shadowedErr)

	default:
		error := errapi.NewServerError("something went wrong")
		httputil.RespondError(w, error)
	}

}
