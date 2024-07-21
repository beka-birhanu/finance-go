package api

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/users"
	commandAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/authentication"
	querieAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_queries/authentication"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/gorilla/mux"
)

type APIServer struct {
	Addr                       string
	UserRepository             persistance.IUserRepository
	UserRegisterCommandHandler commandAuth.IUserRegisterCommandHandler
	UserLoginQueryHandler      querieAuth.IUserLoginQueryHandler
	AuthorizationMiddleware    func(http.Handler) http.Handler
}

func NewAPIServer(addr string, userRepository persistance.IUserRepository, userRegisterCommandHandler commandAuth.IUserRegisterCommandHandler, userQueryHandler querieAuth.IUserLoginQueryHandler, authorizationMiddleware func(http.Handler) http.Handler) *APIServer {
	return &APIServer{
		Addr:                       addr,
		UserRepository:             userRepository,
		UserRegisterCommandHandler: userRegisterCommandHandler,
		UserLoginQueryHandler:      userQueryHandler,
		AuthorizationMiddleware:    authorizationMiddleware,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// Routes that do not require authentication
	publicRouter := router.PathPrefix("/api/v1/public").Subrouter()

	// Routes that require authentication
	protectedRouter := router.PathPrefix("/api/v1").Subrouter()
	protectedRouter.Use(s.AuthorizationMiddleware)

	userHandler := users.NewHandler(s.UserRepository, s.UserRegisterCommandHandler, s.UserLoginQueryHandler)
	userHandler.RegisterPublicRoutes(publicRouter)
	userHandler.RegisterProtectedRoutes(protectedRouter)

	log.Println("Listening on", s.Addr)

	return http.ListenAndServe(s.Addr, router)
}
