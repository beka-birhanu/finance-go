package api

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/users"
	"github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/authentication"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/gorilla/mux"
)

type APIServer struct {
	Addr                      string
	UserRepository            persistance.IUserRepository
	UserRegiserCommandHandler authentication.IUserRegisterCommandHandler
}

func NewAPIServer(addr string, userRepository persistance.IUserRepository, userRegisterCommandHandler authentication.IUserRegisterCommandHandler) *APIServer {
	return &APIServer{
		Addr:                      addr,
		UserRepository:            userRepository,
		UserRegiserCommandHandler: userRegisterCommandHandler,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler := users.NewHandler(s.UserRepository, s.UserRegiserCommandHandler)
	userHandler.RegisterRoutes(subrouter)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.Addr)

	return http.ListenAndServe(s.Addr, router)
}

