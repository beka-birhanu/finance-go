package dto

import (
	"github.com/beka-birhanu/finance-go/application/authentication/common"
)

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func FromAuthResult(authResult *common.AuthResult) *LoginResponse {
	return &LoginResponse{ID: authResult.ID.String(), Username: authResult.Username, Token: authResult.Token}
}
