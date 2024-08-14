// Package user provides HTTP handlers for user-related actions such as registration and login.
// It handles public and protected routes related to user authentication and manages user data.
package user

import (
	"errors"
	"net/http"

	baseapi "github.com/beka-birhanu/finance-go/api/base_handler"
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/user/dto"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	auth "github.com/beka-birhanu/finance-go/application/authentication/common"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	"github.com/gorilla/mux"
)

// Handler manages HTTP requests related to user actions such as registration and login.
type Handler struct {
	baseapi.BaseHandler
	repository      irepository.IUserRepository
	registerHandler icmd.IHandler[*registercmd.Command, *auth.Result]
	loginHandler    iquery.IHandler[*loginqry.Query, *auth.Result]
}

// Config holds the dependencies needed to create a Handler.
type Config struct {
	UserRepository  irepository.IUserRepository
	RegisterHandler icmd.IHandler[*registercmd.Command, *auth.Result]
	LoginHandler    iquery.IHandler[*loginqry.Query, *auth.Result]
}

// NewHandler creates a new Handler with the given configuration.
func NewHandler(config Config) *Handler {
	return &Handler{
		repository:      config.UserRepository,
		registerHandler: config.RegisterHandler,
		loginHandler:    config.LoginHandler,
	}
}

// RegisterPublicRoutes registers public routes for user registration and login.
// These routes are accessible without authentication.
func (h *Handler) RegisterPublic(router *mux.Router) {
	router.HandleFunc("/users/register", h.handleRegistration).Methods(http.MethodPost)
	router.HandleFunc("/users/login", h.handleLogin).Methods(http.MethodPost)
}

// RegisterProtectedRoutes registers routes that require authentication.
// This method is a placeholder and currently does not register any routes.
func (h *Handler) RegisterProtected(router *mux.Router) {}

// handleRegistration processes user registration requests.
// It validates the registration request, creates a registration command,
// and uses the registerHandler to handle the registration logic. On success,
// it sends a response with the authentication result and sets a cookie with the access token.
func (h *Handler) handleRegistration(w http.ResponseWriter, r *http.Request) {
	var registerRequest dto.RegisterRequest
	if err := h.ValidatedBody(r, &registerRequest); err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	registerCommand, err := registercmd.NewCommand(registerRequest.Username, registerRequest.Password)
	if err != nil {
		h.Problem(w, errapi.NewServerError(err.Error()))
		return
	}

	authResult, err := h.registerHandler.Handle(registerCommand)
	if err != nil {
		unwrappedErr := errors.Unwrap(err).(*errdmn.Error)
		switch unwrappedErr.Type() {
		case errdmn.Conflict:
			h.Problem(w, errapi.NewConflict(unwrappedErr.Error()))
		case errdmn.Validation:
			h.Problem(w, errapi.NewBadRequest(unwrappedErr.Error()))
		default:
			h.Problem(w, errapi.NewServerError(unwrappedErr.Error()))
		}
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

	h.RespondWithCookies(w, http.StatusOK, registerResponse, []*http.Cookie{&cookie})
}

// handleLogin processes user login requests.
// It validates the login request, creates a login query, and uses the loginHandler
// to handle the login logic. On success, it sends a response with the authentication result
// and sets a cookie with the access token.
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginUserRequest
	if err := h.ValidatedBody(r, &loginRequest); err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	loginQuery := loginqry.NewQuery(loginRequest.Username, loginRequest.Password)
	authResult, err := h.loginHandler.Handle(loginQuery)
	if err != nil {
		h.Problem(w, errapi.NewAuthentication(err.Error()))
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

	h.RespondWithCookies(w, http.StatusOK, loginResponse, []*http.Cookie{&cookie})
}
