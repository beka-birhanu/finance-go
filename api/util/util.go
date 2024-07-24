package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	apiError "github.com/beka-birhanu/finance-go/api/error"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return apiError.ErrMissingRequestBody
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteError(w http.ResponseWriter, err apiError.Error) {
	WriteJSON(w, int(err.StatusCode), map[string]string{"error": err.Error()})
}

func WriteJSONWithCookie(w http.ResponseWriter, status int, v any, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
	WriteJSON(w, status, v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		}
	}
}

func GetBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func GetIdFromUrl(r *http.Request, paramName string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[paramName]
	if !ok {
		err := apiError.NewErrMissingParam(paramName)
		return uuid.Nil, err
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		err := apiError.NewErrInvalidParam(paramName)
		return uuid.Nil, err
	}
	return id, nil
}

