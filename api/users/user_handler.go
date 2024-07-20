package users

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/users/dto"
	"github.com/beka-birhanu/finance-go/api/utils"
	"github.com/beka-birhanu/finance-go/application/authentication/commands"

	"github.com/beka-birhanu/finance-go/application/authentication/queries"
	commandAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/authentication"
	querieAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_queries/authentication"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	UserRepository             persistance.IUserRepository
	UserRegisterCommandHandler commandAuth.IUserRegisterCommandHandler
	UserQueryHandler           querieAuth.IUserLoginQueryHandler
}

func NewHandler(userRepository persistance.IUserRepository, commandHandler commandAuth.IUserRegisterCommandHandler, queryHandler querieAuth.IUserLoginQueryHandler) *Handler {
	return &Handler{UserRepository: userRepository, UserRegisterCommandHandler: commandHandler, UserQueryHandler: queryHandler}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(
		"/users/register",
		h.HandleUserRegisteration,
	).Methods(http.MethodPost)
}

func (h *Handler) HandleUserRegisteration(w http.ResponseWriter, r *http.Request) {
	var registerRequest dto.RegisterRequest

	if err := utils.ParseJSON(r, &registerRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(registerRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	registerCommand, err := commands.NewUserRegisterCommand(registerRequest.Username, registerRequest.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.UserRegisterCommandHandler.Handle(registerCommand); err != nil {
		if err == commands.ErrUsernameInUse {
			utils.WriteError(w, http.StatusConflict, err)
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"token": "token"})
}
