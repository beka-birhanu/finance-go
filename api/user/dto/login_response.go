package dto

import (
	"github.com/beka-birhanu/finance-go/application/authentication/common"
)

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// FromAuthResult extracts the info for the login response from the given
// auth.Result and map them to new LoginResponse
func FromAuthResult(authResult *auth.Result) *LoginResponse {
	return &LoginResponse{ID: authResult.ID.String(), Username: authResult.Username}
}
