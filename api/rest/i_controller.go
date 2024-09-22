// Package api defines an interface for controllers to register routes
// with different access levels.
package api

import "github.com/gorilla/mux"

// IController outlines methods for route registration:
// - Public: No authentication required.
// - Protected: Requires authentication.
type IController interface {
	// RegisterPublic sets up public routes.
	RegisterPublic(router *mux.Router)

	// RegisterProtected sets up routes that require authentication.
	RegisterProtected(router *mux.Router)
}
