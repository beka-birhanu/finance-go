package commands

import (
	"fmt"

	"github.com/google/uuid"
)

type UserRegisterCommand struct {
	ID       uuid.UUID
	Username string
	Password string
}

func NewUserRegisterCommand(username, password string) (*UserRegisterCommand, error) {
	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("server error")
	}

	return &UserRegisterCommand{ID: id, Username: username, Password: password}, nil
}
