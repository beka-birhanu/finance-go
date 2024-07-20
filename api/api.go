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
}

func NewAPIServer(addr string, userRepository persistance.IUserRepository, userRegisterCommandHandler commandAuth.IUserRegisterCommandHandler, userQueryHandler querieAuth.IUserLoginQueryHandler) *APIServer {
	return &APIServer{
		Addr:                       addr,
		UserRepository:             userRepository,
		UserRegisterCommandHandler: userRegisterCommandHandler,
		UserLoginQueryHandler:      userQueryHandler,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler := users.NewHandler(s.UserRepository, s.UserRegisterCommandHandler, s.UserLoginQueryHandler)
	userHandler.RegisterRoutes(subrouter)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.Addr)

	return http.ListenAndServe(s.Addr, router)
}
