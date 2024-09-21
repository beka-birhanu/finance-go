// Package baseapi provides a base handler for HTTP requests, including
// functionalities for request validation, error handling, and response management.
package baseapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/beka-birhanu/finance-go/api/middleware"
	"github.com/dgrijalva/jwt-go"
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Validate is a global validator instance used for validating structs.
var Validate = validator.New()

// BaseHandler is a base struct for all HTTP request handlers, providing
// basic functionalities such as request validation, error handling, and response formatting.
type BaseHandler struct{}

// ValidatedBody validates the body of an HTTP request and populates the provided
// interface. It returns an error if the request body is invalid or fails validation.
func (h *BaseHandler) ValidatedBody(r *http.Request, s interface{}) error {
	if err := h.ParseJSON(r, s); err != nil {
		return err
	}

	if err := Validate.Struct(s); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return errapi.NewBadRequest(fmt.Sprintf("invalid payload: %v", validationErrs))
		}
		return errapi.NewServerError("request body validation failed")
	}

	return nil
}

// Problem handles errors by writing an appropriate response to the HTTP response writer.
// It masks certain internal error details to avoid exposing them to the client.
func (h *BaseHandler) Problem(w http.ResponseWriter, err errapi.Error) {
	var shadowedErr errapi.Error
	switch err.StatusCode() {
	case errapi.BadRequest, errapi.Conflict, errapi.NotFound, errapi.Forbidden:
		shadowedErr = err
	case errapi.Authentication:
		shadowedErr = errapi.NewAuthentication("invalid credentials")
	default:
		shadowedErr = errapi.NewServerError("something went wrong")
	}
	h.RespondError(w, shadowedErr)
}

// ParseJSON decodes the JSON body of an HTTP request into the provided interface.
// It returns an error if the request body is missing or invalid.
func (h *BaseHandler) ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errapi.NewBadRequest("request body is missing")
	}

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errapi.NewBadRequest("invalid request body")
	}
	return nil
}

// RespondError writes an error response with the appropriate status code and error message
// to the HTTP response writer.
func (h *BaseHandler) RespondError(w http.ResponseWriter, err errapi.Error) {
	h.Respond(w, err.StatusCode(), map[string]string{"error": err.Error()})
}

// RespondWithCookies writes a response with cookies to the HTTP response writer.
// It adds the provided cookies to the response.
func (h *BaseHandler) RespondWithCookies(w http.ResponseWriter, status int, v interface{}, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
	h.Respond(w, status, v)
}

// RespondWithLocation writes a response with a Location header to the HTTP response writer.
// It sets the Location header to the provided resource location.
func (h *BaseHandler) RespondWithLocation(w http.ResponseWriter, status int, v interface{}, resourceLocation string) {
	w.Header().Set("Location", resourceLocation)
	h.Respond(w, status, v)
}

// Respond writes a JSON response to the HTTP response writer with the given status code.
// It encodes the provided interface as JSON in the response body.
func (h *BaseHandler) Respond(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			http.Error(w, "failed to encode JSON response", http.StatusInternalServerError)
		}
	}
}

// BaseURL returns the base URL of the request, including the scheme and host.
func (h *BaseHandler) BaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

// UUIDParam retrieves a UUID path parameter from the request URL. It returns an error
// if the parameter is missing or not a valid UUID.
func (h *BaseHandler) UUIDParam(r *http.Request, paramName string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[paramName]
	if !ok {
		return uuid.Nil, errapi.NewBadRequest(fmt.Sprintf("path parameter %v is missing", paramName))
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errapi.NewBadRequest(fmt.Sprintf("path parameter %v is of invalid format", paramName))
	}
	return id, nil
}

// MatchPathUserIdctxUserId checks if the user ID from the request context matches the provided user ID.
// It returns an error if the IDs do not match or if there is an issue retrieving the user ID from the context.
func (h *BaseHandler) MatchPathUserIdctxUserId(r *http.Request, pathId uuid.UUID) error {
	claims, ok := r.Context().Value(middleware.ContextUserClaims).(jwt.MapClaims)
	if !ok {
		return errapi.NewServerError("error on retrieving user id from context")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok || pathId.String() != userIDStr {
		return errapi.NewForbidden("The response does not belong to the user requesting.")
	}

	return nil
}

// StringQueryParam retrieves a string query parameter from the request URL.
func (h *BaseHandler) StringQueryParam(r *http.Request, paramName string) string {
	return r.URL.Query().Get(paramName)
}

// IntQueryParam retrieves an integer query parameter from the request URL.
// It returns an error if the parameter is not a valid integer.
func (h *BaseHandler) IntQueryParam(r *http.Request, paramName string) (int, error) {
	param := r.URL.Query().Get(paramName)
	if param == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return 0, errapi.NewBadRequest(fmt.Sprintf("invalid query parameter %s: %v", paramName, err))
	}
	return val, nil
}
