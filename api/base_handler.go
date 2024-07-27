package api

import (
	"errors"
	"fmt"
	"net/http"

	errapi "github.com/beka-birhanu/finance-go/api/error"
	httputil "github.com/beka-birhanu/finance-go/api/http_util"
	"github.com/go-playground/validator/v10"
)

type BaseHandler struct{}

// ValidatedBody validates the body of an HTTP request.
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
		return errapi.NewServerError("validation failed")
	}

	return nil
}

