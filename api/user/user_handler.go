package user

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/finance-go/api"
	errapi "github.com/beka-birhanu/finance-go/api/error"
	httputil "github.com/beka-birhanu/finance-go/api/http_util"
	"github.com/beka-birhanu/finance-go/api/user/dto"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	auth "github.com/beka-birhanu/finance-go/application/authentication/common"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	api.BaseHandler
	repository      irepository.IUserRepository
	registerHandler icmd.IHandler[*registercmd.Command, *auth.Result]
	loginHandler    iquery.IHandler[*loginqry.Query, *auth.Result]
}

type Config struct {
	UserRepository  irepository.IUserRepository
	RegisterHandler icmd.IHandler[*registercmd.Command, *auth.Result]
	LoginHandler    icmd.IHandler[*loginqry.Query, *auth.Result]
}

func NewHandler(config Config) *Handler {
	return &Handler{
		repository:      config.UserRepository,
		registerHandler: config.RegisterHandler,
		loginHandler:    config.LoginHandler,
	}
}

func (h *Handler) RegisterPublicRoutes(router *mux.Router) {
	router.HandleFunc(
		"/users/register",
		h.handleRegistration,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/users/login",
		h.handleLogin,
	).Methods(http.MethodPost)
}

func (h *Handler) RegisterProtectedRoutes(router *mux.Router) {}

func (h *Handler) handleRegistration(w http.ResponseWriter, r *http.Request) {
	var registerRequest dto.RegisterRequest

	if err := httputil.ParseJSON(r, &registerRequest); err != nil {
		httputil.RespondError(w, err.(errapi.Error))
		return
	}

	if err := httputil.Validate.Struct(registerRequest); err != nil {
		errors := err.(validator.ValidationErrors)

		errResponse := errapi.NewBadRequest(fmt.Sprintf("invalid payload: %v", errors))
		httputil.RespondError(w, errResponse)
		return
	}

	registerCommand, err := registercmd.NewCommand(registerRequest.Username, registerRequest.Password)
	if err != nil {
		apiErr := errapi.NewServerError("unexpected server error")
		httputil.RespondError(w, apiErr)
		return
	}

	authResult, err := h.registerHandler.Handle(registerCommand)
	if err != nil {
		err := errapi.NewBadRequest(err.Error())
		httputil.RespondError(w, err)
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

	httputil.RespondWithCookies(w, http.StatusOK, registerResponse, []*http.Cookie{&cookie})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginUserRequest

	if err := httputil.ParseJSON(r, &loginRequest); err != nil {
		httputil.RespondError(w, err.(errapi.Error))
		return
	}

	if err := httputil.Validate.Struct(loginRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		err := errapi.NewBadRequest(fmt.Sprintf("invalid payload: %v", errors))
		httputil.RespondError(w, err)
		return
	}

	loginQuery := loginqry.NewQuery(loginRequest.Username, loginRequest.Password)

	authResult, err := h.loginHandler.Handle(loginQuery)
	if err != nil {
		apiErr := errapi.NewBadRequest(err.Error())
		httputil.RespondError(w, apiErr)
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

	httputil.RespondWithCookies(w, http.StatusOK, loginResponse, []*http.Cookie{&cookie})
}
