package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/users"
	"github.com/beka-birhanu/finance-go/infrastructure/repositories"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userRepository := repositories.NewUserRepository(s.db)
	userHandler := users.NewHandler(userRepository)
	userHandler.RegisterRoutes(subrouter)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
