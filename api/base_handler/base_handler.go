package baseapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Validate is a global validator instance for validating structs.
var Validate = validator.New()

// BaseHandler is a base struct for all HTTP request handlers that provides basic HTTP functionalities.
type BaseHandler struct{}

// ValidatedBody validates the body of an HTTP request and populates the given interface.
func (h *BaseHandler) ValidatedBody(r *http.Request, s interface{}) error {
	// Parse JSON body
	if err := h.ParseJSON(r, s); err != nil {
		var apiErr errapi.Error
		if errors.As(err, &apiErr) {
			return apiErr
		}
		return errapi.NewServerError("unable to parse request body")
	}

	// Validate the parsed struct
	if err := Validate.Struct(s); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return errapi.NewBadRequest(fmt.Sprintf("invalid payload: %v", validationErrs))
		}
		return errapi.NewServerError("request body validation failed")
	}

	return nil
}

// Problem handles errors by writing an appropriate response to the HTTP response writer.
func (h *BaseHandler) Problem(w http.ResponseWriter, err errapi.Error) {
	// Shadow error messages if they are not supposed to be seen.
	switch err.StatusCode() {
	case errapi.BadRequest:
		h.RespondError(w, err)
	case errapi.Conflict:
		h.RespondError(w, err)
	case errapi.NotFound:
		h.RespondError(w, err)
	case errapi.Forbidden:
		h.RespondError(w, err)
	// Client should not know what caused the failure.
	case errapi.Authentication:
		shadowedErr := errapi.NewAuthentication("invalid credentials")
		h.RespondError(w, shadowedErr)
	// Client should not know what caused the failure.
	case errapi.ServerError:
		shadowedErr := errapi.NewServerError("something went wrong")
		h.RespondError(w, shadowedErr)
	default:
		error := errapi.NewServerError("something went wrong")
		h.RespondError(w, error)
	}
}

// ParseJSON decodes the JSON body of an HTTP request into the provided interface.
// It returns a BadRequest error if the request body is missing or invalid.
func (h *BaseHandler) ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errapi.NewBadRequest("request body is missing")
	}

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return errapi.NewBadRequest("invalid request body")
	}
	return nil
}

// RespondError writes an error response to the HTTP response writer.
// The response contains the error status code and message.
func (h *BaseHandler) RespondError(w http.ResponseWriter, err errapi.Error) {
	h.Respond(w, err.StatusCode(), map[string]string{"error": err.Error()})
}

// RespondWithCookies writes a response with cookies to the HTTP response writer.
// It sets the provided cookies in the response header before writing the response.
func (h *BaseHandler) RespondWithCookies(w http.ResponseWriter, status int, v interface{}, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
	h.Respond(w, status, v)
}

// RespondWithLocation writes a response with a Location header to the HTTP response writer.
// It sets the Location header before writing the response.
func (h *BaseHandler) RespondWithLocation(w http.ResponseWriter, status int, v interface{}, resourceLocation string) {
	w.Header().Set("Location", resourceLocation)
	h.Respond(w, status, v)
}

// Respond writes a JSON response to the HTTP response writer.
// It sets the Content-Type header to application/json and writes the provided status code and value.
func (h *BaseHandler) Respond(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		err := json.NewEncoder(w).Encode(v)
		if err != nil {
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

// UUIDParam retrieves a UUID path parameter from the request URL.
// It returns a BadRequest error if the parameter is missing or invalid.
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

// MatchPathUserIdctxUserId returns err if the userId for a request's context and provided are not the same.
func (h *BaseHandler) MatchPathUserIdctxUserId(r *http.Request, pathId uuid.UUID) error {
	// Extract userId for context and match with the userId from URL.
	ctx := r.Context()
	claims, ok := ctx.Value(middleware.ContextUserClaims).(jwt.MapClaims)
	if !ok {
		err := errapi.NewServerError("error on retrieving user id from context")
		return err
	}

	// Accessing the user_id as string and then parse it to uuid.UUID
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		err := errapi.NewForbidden("The response does not belong to the user requesting.")
		return err
	}

	// Parse the userIDStr to uuid.UUID
	ctxUserId, err := uuid.Parse(userIDStr)
	if err != nil {
		err := errapi.NewServerError("invalid user id format")
		return err
	}

	if ctxUserId != pathId {
		err := errapi.NewForbidden("The response does not belong to the user requesting.")
		return err
	}

	return nil
}

func (h *BaseHandler) StringQueryParam(r *http.Request, paramName string) string {
	return r.URL.Query().Get(paramName)
}

func (h *BaseHandler) IntQueryParam(r *http.Request, paramName string) (int, error) {
	param := r.URL.Query().Get(paramName)
	if param == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		err = errapi.NewBadRequest(fmt.Sprintf("invalid query parameter %s: %v", paramName, err))
	}

	return val, err
}
