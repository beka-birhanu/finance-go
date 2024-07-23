package commands

type UserRegisterCommand struct {
	Username string
	Password string
}

func NewUserRegisterCommand(username, password string) (*UserRegisterCommand, error) {
	return &UserRegisterCommand{Username: username, Password: password}, nil
}
