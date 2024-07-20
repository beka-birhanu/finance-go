package common

import "github.com/google/uuid"

type AuthResult struct {
	ID       uuid.UUID
	Username string
	Token    string
}

func NewAuthResult(id uuid.UUID, username, token string) *AuthResult {
	return &AuthResult{
		ID:       id,
		Username: username,
		Token:    token,
	}
}
