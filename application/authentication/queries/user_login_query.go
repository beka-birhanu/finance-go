package queries

type UserLoginQuery struct {
	Username string
	Password string
}

func NewUserLoginQuery(username, password string) *UserLoginQuery {
	return &UserLoginQuery{Username: username, Password: password}
}
