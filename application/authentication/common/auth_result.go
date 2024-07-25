package auth

import "github.com/google/uuid"

// Result represents the result of a successful authentication.
type Result struct {
	// ID is the unique identifier of the authenticated user.
	ID uuid.UUID

	// Username is the username of the authenticated user.
	Username string

	// Token is the authentication token issued to the authenticated user.
	Token string
}

// NewResult creates and return a new Result instance with the provided
// user ID, username, and token.
func NewResult(id uuid.UUID, username, token string) *Result {
	return &Result{
		ID:       id,
		Username: username,
		Token:    token,
	}
}

