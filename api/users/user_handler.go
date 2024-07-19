package users

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/infrastructure/repositories"
	"github.com/gorilla/mux"
)

type Handler struct {
	UserRepository *repositories.UserRepository
}

func NewHandler(userRepository *repositories.UserRepository) *Handler {
	return &Handler{UserRepository: userRepository}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(
		"/users/register",
		h.handleUserRegisteration,
	).Methods(http.MethodPost)
}

func (h *Handler) handleUserRegisteration(w http.ResponseWriter, r *http.Request) {
	log.Println("working")
}
