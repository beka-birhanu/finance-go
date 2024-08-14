package registercmd

// Command hold feilds required for registering users.
type Command struct {
	// Username is username of the user to be registerd
	Username string

	// Password is the plain text password of the user to be registerd
	Password string
}

// NewCommand returns a new command for user registeration.
func NewCommand(username, password string) (*Command, error) {
	return &Command{Username: username, Password: password}, nil
}
