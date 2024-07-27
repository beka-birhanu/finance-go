package user

import (
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
	"github.com/gorilla/mux"
)

// Handler manages user-related HTTP requests.
type Handler struct {
	api.BaseHandler
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

// RegisterPublicRoutes registers the public routes for user-related actions.
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

// RegisterProtectedRoutes registers the protected routes for user-related actions.
func (h *Handler) RegisterProtectedRoutes(router *mux.Router) {}

func (h *Handler) handleRegistration(w http.ResponseWriter, r *http.Request) {
	var registerRequest dto.RegisterRequest

	// Populate registerRequest from request body
	if err := h.ValidatedBody(r, &registerRequest); err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	registerCommand, err := registercmd.NewCommand(registerRequest.Username, registerRequest.Password)
	if err != nil {
		apiErr := errapi.NewServerError(err.Error())
		h.Problem(w, apiErr)
		return
	}

	authResult, err := h.registerHandler.Handle(registerCommand)
	if err != nil {
		err := errapi.NewBadRequest(err.Error())
		h.Problem(w, err)
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

	// Populate loginRequest from request body
	if err := h.ValidatedBody(r, &loginRequest); err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	loginQuery := loginqry.NewQuery(loginRequest.Username, loginRequest.Password)

	authResult, err := h.loginHandler.Handle(loginQuery)
	if err != nil {
		apiErr := errapi.NewAuthentication(err.Error())
		h.Problem(w, apiErr)
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

