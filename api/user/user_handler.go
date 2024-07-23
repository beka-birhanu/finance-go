package users

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/users/dto"
	"github.com/beka-birhanu/finance-go/api/utils"
	"github.com/beka-birhanu/finance-go/application/authentication/commands"
	"github.com/beka-birhanu/finance-go/domain/domain_errors"

	"github.com/beka-birhanu/finance-go/application/authentication/queries"
	commandAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/authentication"
	querieAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_queries/authentication"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UsersHandler struct {
	userRepository             persistance.IUserRepository
	userRegisterCommandHandler commandAuth.IUserRegisterCommandHandler
	userQueryHandler           querieAuth.IUserLoginQueryHandler
}

func NewHandler(userRepository persistance.IUserRepository, commandHandler commandAuth.IUserRegisterCommandHandler, queryHandler querieAuth.IUserLoginQueryHandler) *UsersHandler {
	return &UsersHandler{userRepository: userRepository, userRegisterCommandHandler: commandHandler, userQueryHandler: queryHandler}
}

func (h *UsersHandler) RegisterPublicRoutes(router *mux.Router) {
	router.HandleFunc(
		"/users/register",
		h.handleUserRegistration,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/users/login",
		h.handleUserLogin,
	).Methods(http.MethodPost)
}

func (h *UsersHandler) RegisterProtectedRoutes(router *mux.Router) {}

func (h *UsersHandler) handleUserRegistration(w http.ResponseWriter, r *http.Request) {
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

	authResult, err := h.userRegisterCommandHandler.Handle(registerCommand)
	if err != nil {
		switch err {
		case domain_errors.ErrUsernameConflict:
			utils.WriteError(w, http.StatusConflict, err)
		case domain_errors.ErrWeakPassword,
			domain_errors.ErrUsernameTooLong,
			domain_errors.ErrUsernameTooShort,
			domain_errors.ErrUsernameInvalidFormat:
			utils.WriteError(w, http.StatusBadRequest, err)
		default:
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	registorResponse := dto.FromAuthResult(authResult)
	cookie := http.Cookie{
		Name:     "accessToken",
		Value:    authResult.Token,
		Path:     "/",
		MaxAge:   24 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	utils.WriteJSONWithCookie(w, http.StatusOK, registorResponse, []*http.Cookie{&cookie})
}

func (h *UsersHandler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginUserRequest

	if err := utils.ParseJSON(r, &loginRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(loginRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	loginQuery := queries.NewUserLoginQuery(loginRequest.Username, loginRequest.Password)

	authResult, err := h.userQueryHandler.Handle(loginQuery)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	loginResponse := dto.FromAuthResult(authResult)

	cookie := http.Cookie{
		Name:     "accessToken",
		Value:    authResult.Token,
		Path:     "/",
		MaxAge:   24 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	utils.WriteJSONWithCookie(w, http.StatusOK, loginResponse, []*http.Cookie{&cookie})
}
