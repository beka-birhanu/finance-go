package user

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/user/dto"
	"github.com/beka-birhanu/finance-go/api/util"
	"github.com/beka-birhanu/finance-go/application/authentication/command"
	authCommand "github.com/beka-birhanu/finance-go/application/authentication/command"
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/authentication/query"
	authQuery "github.com/beka-birhanu/finance-go/application/authentication/query"
	handlerInterface "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	domainError "github.com/beka-birhanu/finance-go/domain/error"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UsersHandler struct {
	userRepository             repository.IUserRepository
	userRegisterCommandHandler handlerInterface.ICommandHandler[*authCommand.UserRegisterCommand, *common.AuthResult]
	userLoginQueryHandler      handlerInterface.ICommandHandler[*authQuery.UserLoginQuery, *common.AuthResult]
}

func NewHandler(
	userRepository repository.IUserRepository,
	userRegisterCommandHandler handlerInterface.ICommandHandler[*authCommand.UserRegisterCommand,
		*common.AuthResult],
	userLoginQueryHandler handlerInterface.ICommandHandler[*authQuery.UserLoginQuery, *common.AuthResult],
) *UsersHandler {
	return &UsersHandler{
		userRepository:             userRepository,
		userRegisterCommandHandler: userRegisterCommandHandler,
		userLoginQueryHandler:      userLoginQueryHandler,
	}
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

	if err := util.ParseJSON(r, &registerRequest); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := util.Validate.Struct(registerRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	registerCommand, err := command.NewUserRegisterCommand(registerRequest.Username, registerRequest.Password)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	authResult, err := h.userRegisterCommandHandler.Handle(registerCommand)
	if err != nil {
		switch err {
		case domainError.ErrUsernameConflict:
			util.WriteError(w, http.StatusConflict, err)
		case domainError.ErrWeakPassword,
			domainError.ErrUsernameTooLong,
			domainError.ErrUsernameTooShort,
			domainError.ErrUsernameInvalidFormat:
			util.WriteError(w, http.StatusBadRequest, err)
		default:
			util.WriteError(w, http.StatusInternalServerError, err)
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

	util.WriteJSONWithCookie(w, http.StatusOK, registorResponse, []*http.Cookie{&cookie})
}

func (h *UsersHandler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginUserRequest

	if err := util.ParseJSON(r, &loginRequest); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := util.Validate.Struct(loginRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	loginQuery := query.NewUserLoginQuery(loginRequest.Username, loginRequest.Password)

	authResult, err := h.userLoginQueryHandler.Handle(loginQuery)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
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

	util.WriteJSONWithCookie(w, http.StatusOK, loginResponse, []*http.Cookie{&cookie})
}
