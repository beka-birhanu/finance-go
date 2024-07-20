package api

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/users"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/gorilla/mux"
)

type APIServer struct {
	Addr           string
	UserRepository persistance.IUserRepository
}

func NewAPIServer(addr string, userRepository persistance.IUserRepository) *APIServer {
	return &APIServer{
		Addr:           addr,
		UserRepository: userRepository,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler := users.NewHandler(s.UserRepository)
	userHandler.RegisterRoutes(subrouter)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.Addr)

	return http.ListenAndServe(s.Addr, router)
}

