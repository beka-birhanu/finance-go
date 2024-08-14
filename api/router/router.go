// Package router provides functionality to set up and run the HTTP server,
// manage routes, and apply middleware based on access levels.
//
// It configures and initializes routes with varying access requirements:
// - Public routes: Accessible without authentication.
// - Protected routes: Require authentication.
// - Privileged routes: Require both authentication and admin privileges.
package router

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/gorilla/mux"
)

// Router manages the HTTP server and its dependencies,
// including controllers and JWT authentication.
type Router struct {
	addr                    string
	baseURL                 string
	controllers             []api.IController
	authorizationMiddleware func(http.Handler) http.Handler
}

// Config holds configuration settings for creating a new Router instance.
type Config struct {
	Addr                    string            // Address to listen on
	BaseURL                 string            // Base URL for API routes
	Controllers             []api.IController // List of controllers
	AuthorizationMiddleware func(http.Handler) http.Handler
}

// NewRouter creates a new Router instance with the given configuration.
// It initializes the router with address, base URL, controllers, and JWT service.
func NewRouter(config Config) *Router {
	return &Router{
		addr:                    config.Addr,
		baseURL:                 config.BaseURL,
		controllers:             config.Controllers,
		authorizationMiddleware: config.AuthorizationMiddleware,
	}
}

// Run starts the HTTP server and sets up routes with different access levels.
//
// Routes are grouped and managed under the base URL, with the following access levels:
// - Public routes: No authentication required.
// - Protected routes: Authentication required.
// - Privileged routes: Authentication and admin privileges required.
func (r *Router) Run() error {
	router := mux.NewRouter()

	// Setting up routes under baseURL
	api := router.PathPrefix("/api").Subrouter()

	{
		// Public routes (accessible without authentication)
		publicRoutes := api.PathPrefix("/v1").Subrouter()
		{
			for _, c := range r.controllers {
				c.RegisterPublic(publicRoutes)
			}
		}

		// Protected routes (authentication required)
		protectedRoutes := api.PathPrefix("/v1").Subrouter()
		protectedRoutes.Use(r.authorizationMiddleware)
		{
			for _, c := range r.controllers {
				c.RegisterProtected(protectedRoutes)
			}
		}

	}

	log.Println("Listening on", r.addr)
	return http.ListenAndServe(r.addr, router)
}
