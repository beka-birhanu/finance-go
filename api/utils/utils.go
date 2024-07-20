package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
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

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
