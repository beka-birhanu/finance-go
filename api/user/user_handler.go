package user

import (
	"fmt"
	"net/http"

	apiError "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/user/dto"
	"github.com/beka-birhanu/finance-go/api/util"
	"github.com/beka-birhanu/finance-go/application/authentication/command"
	authCommand "github.com/beka-birhanu/finance-go/application/authentication/command"
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	"github.com/beka-birhanu/finance-go/application/authentication/query"
	authQuery "github.com/beka-birhanu/finance-go/application/authentication/query"
	handlerInterface "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
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
		apiErr := apiError.NewErrBadRequest(err.Error())
		util.WriteError(w, apiErr)
		return
	}

	if err := util.Validate.Struct(registerRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		apiErr := apiError.NewErrValidation(fmt.Sprintf("invalid payload: %v", errors))
		util.WriteError(w, apiErr)
		return
	}

	registerCommand, err := command.NewUserRegisterCommand(registerRequest.Username, registerRequest.Password)
	if err != nil {
		apiErr := apiError.NewErrServer(err.Error())
		util.WriteError(w, apiErr)
		return
	}

	authResult, err := h.userRegisterCommandHandler.Handle(registerCommand)
	if err != nil {
		err := apiError.ErrToAPIError(err)
		util.WriteError(w, err)
		return
	}

	registerResponse := dto.FromAuthResult(authResult)
	cookie := http.Cookie{
		Name:     "accessToken",
		Value:    authResult.Token,
		Path:     "/",
		MaxAge:   24 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	util.WriteJSONWithCookie(w, http.StatusOK, registerResponse, []*http.Cookie{&cookie})
}

func (h *UsersHandler) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginUserRequest

	if err := util.ParseJSON(r, &loginRequest); err != nil {
		apiErr := apiError.ErrToAPIError(err)
		util.WriteError(w, apiErr)
		return
	}

	if err := util.Validate.Struct(loginRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		apiErr := apiError.NewErrValidation(fmt.Sprintf("invalid payload: %v", errors))
		util.WriteError(w, apiErr)
		return
	}

	loginQuery := query.NewUserLoginQuery(loginRequest.Username, loginRequest.Password)

	authResult, err := h.userLoginQueryHandler.Handle(loginQuery)
	if err != nil {
		apiErr := apiError.ErrToAPIError(err)
		util.WriteError(w, apiErr)
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

