package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"

	apiError "github.com/beka-birhanu/finance-go/api/error"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Validate is a global validator instance for validating structs.
var Validate = validator.New()

// ParseJSON decodes the JSON body of an HTTP request into the provided interface.
// It returns a BadRequest error if the request body is missing.
func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return apiError.NewBadRequest("Request body is missing")
	}

	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return apiError.NewBadRequest("Invalid request body")
	}
	return nil
}

// RespondError writes an error response to the HTTP response writer only from apiError.
// The response contains the error status code and message.
func RespondError(w http.ResponseWriter, err apiError.Error) {
	Respond(w, err.StatusCode(), map[string]string{"error": err.Error()})
}

// RespondWithCookies writes a response with cookies to the HTTP response writer.
// It sets the provided cookies in the response header before writing the response.
func RespondWithCookies(w http.ResponseWriter, status int, v any, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
	Respond(w, status, v)
}

// RespondWithLocation writes a response with a Location header to the HTTP response writer.
func ResondWithLocation(w http.ResponseWriter, status int, v any, resourceLocation string) {
	w.Header().Set("Location", resourceLocation)
	Respond(w, status, v)
}

// Respond writes a JSON response to the HTTP response writer.
// It sets the Content-Type header to application/json and writes the provided status code and value.
func Respond(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		err := json.NewEncoder(w).Encode(v)
		if err != nil {
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		}
	}
}

// BaseURL returns the base URL of the request, including the scheme and host.
func BaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

// UUIDParam retrieves a UUID path parameter from the request URL.
// It returns a BadRequest error if the parameter is missing or invalid.
func UUIDParam(r *http.Request, paramName string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[paramName]
	if !ok {
		return uuid.Nil, apiError.NewBadRequest(fmt.Sprintf("Path parameter %v is missing", paramName))
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, apiError.NewBadRequest(fmt.Sprintf("Path parameter %v is of invalid format", paramName))
	}
	return id, nil
}
