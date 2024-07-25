// Package loginqry provides functionality for handling login queries.
package loginqry

// Query represents a login query containing a username and password.
type Query struct {
	// Username is the username provided for the login attempt.
	Username string

	// Password is the password provided for the login attempt.
	Password string
}

// NewQuery creates and return a new Query instance with the provided
// username and password.
func NewQuery(username, password string) *Query {
	return &Query{Username: username, Password: password}
}
