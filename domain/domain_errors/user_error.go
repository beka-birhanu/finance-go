package domain_errors

import (
	"errors"
)

var (
	UsernameConflict = errors.New("username already taken")
)
