package auth

import "github.com/google/uuid"

// AuthResult represents the result of a successful authentication.
type AuthResult struct {
	// ID is the unique identifier of the authenticated user.
	ID uuid.UUID

	// Username is the username of the authenticated user.
	Username string

	// Token is the authentication token issued to the authenticated user.
	Token string
}

// Result creates and return a new AuthResult instance with the provided
// user ID, username, and token.
func Result(id uuid.UUID, username, token string) *AuthResult {
	return &AuthResult{
		ID:       id,
		Username: username,
		Token:    token,
	}
}

