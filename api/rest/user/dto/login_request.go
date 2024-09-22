package dto

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
