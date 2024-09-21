// Package router provides functionality to set up and run the HTTP server,
// manage routes, and apply middleware based on access levels.
//
// It configures and initializes routes with varying access requirements:
// - Public routes: Accessible without authentication.
// - Protected routes: Require authentication.
package router

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	api "github.com/beka-birhanu/finance-go/api/rest"
	"github.com/gorilla/mux"
)

// Router manages the HTTP server and its dependencies,
// including controllers and JWT authentication.
type Router struct {
	addr                    string
	baseURL                 string
	restfullControllers     []api.IController
	graphQlController       http.Handler
	authorizationMiddleware func(http.Handler) http.Handler
	rateLimitMiddleware     func(http.Handler) http.Handler
}

// Config holds configuration settings for creating a new Router instance.
type Config struct {
	Addr                    string            // Address to listen on
	BaseURL                 string            // Base URL for API routes
	RestfullControllers     []api.IController // List of controllers
	GraphQlController       http.Handler
	AuthorizationMiddleware func(http.Handler) http.Handler
	RateLimitMiddleware     func(http.Handler) http.Handler
}

// NewRouter creates a new Router instance with the given configuration.
// It initializes the router with address, base URL, controllers, and JWT service.
func NewRouter(config Config) *Router {
	return &Router{
		addr:                    config.Addr,
		baseURL:                 config.BaseURL,
		restfullControllers:     config.RestfullControllers,
		graphQlController:       config.GraphQlController,
		authorizationMiddleware: config.AuthorizationMiddleware,
		rateLimitMiddleware:     config.RateLimitMiddleware,
	}
}

// Run starts the HTTP server and sets up routes with different access levels.
//
// Routes are grouped and managed under the base URL, with the following access levels:
// - Public routes: No authentication required.
// - Protected routes: Authentication required.
func (r *Router) Run() error {
	router := mux.NewRouter()
	router.Use((r.rateLimitMiddleware))

	// Setting up routes under baseURL
	api := router.PathPrefix("/api").Subrouter()

	{
		// Public routes (accessible without authentication)
		publicRoutes := api.PathPrefix("/v1").Subrouter()
		{
			for _, c := range r.restfullControllers {
				c.RegisterPublic(publicRoutes)
			}
		}

		// Protected routes (authentication required)
		protectedRoutes := api.PathPrefix("/v1").Subrouter()
		protectedRoutes.Use(r.authorizationMiddleware)
		{
			for _, c := range r.restfullControllers {
				c.RegisterProtected(protectedRoutes)
			}
		}

	}
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", r.graphQlController)

	log.Println("Listening on", r.addr)
	return http.ListenAndServe(r.addr, router)
}
