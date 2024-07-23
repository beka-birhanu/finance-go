package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteJSONWithCookie(w http.ResponseWriter, status int, v any, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
	WriteJSON(w, http.StatusOK, v)
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
	// Extract user ID from URL path
	vars := mux.Vars(r)
	idStr, ok := vars[paramName]
	if !ok {
		return uuid.Nil, fmt.Errorf("missing user %s in path", paramName)
	}

	// Parse user ID to UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user %s format: %v", paramName, err)
	}
	return id, nil
}
