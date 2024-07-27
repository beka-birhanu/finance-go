package api

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/expense"
	"github.com/beka-birhanu/finance-go/api/user"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	auth "github.com/beka-birhanu/finance-go/application/authentication/common"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/gorilla/mux"
)

// APIServer is the struct for managing the HTTP server and its dependencies.
type APIServer struct {
	addr                     string
	userRepository           irepository.IUserRepository
	userRegisterHandler      icmd.IHandler[*registercmd.Command, *auth.Result]
	userLoginQueryHandler    iquery.IHandler[*loginqry.Query, *auth.Result]
	authorizationMiddleware  func(http.Handler) http.Handler
	addExpenseCommandHandler icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	timeService              itimeservice.IService
}

// Config is the struct for configuring the APIServer.
type Config struct {
	Addr                     string
	UserRepository           irepository.IUserRepository
	UserRegisterHandler      icmd.IHandler[*registercmd.Command, *auth.Result]
	UserLoginQueryHandler    iquery.IHandler[*loginqry.Query, *auth.Result]
	AuthorizationMiddleware  func(http.Handler) http.Handler
	AddExpenseCommandHandler icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	TimeService              itimeservice.IService
}

// NewAPIServer creates a new instance of APIServer with the given configuration.
func NewAPIServer(config Config) *APIServer {
	return &APIServer{
		addr:                     config.Addr,
		userRepository:           config.UserRepository,
		userRegisterHandler:      config.UserRegisterHandler,
		userLoginQueryHandler:    config.UserLoginQueryHandler,
		authorizationMiddleware:  config.AuthorizationMiddleware,
		addExpenseCommandHandler: config.AddExpenseCommandHandler,
		timeService:              config.TimeService,
	}
}

// Run starts the HTTP server and sets up the routes.
func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// Public routes that do not require authentication
	publicRouter := router.PathPrefix("/api/v1/public").Subrouter()

	// Protected routes that require authentication
	protectedRouter := router.PathPrefix("/api/v1").Subrouter()
	protectedRouter.Use(s.authorizationMiddleware)

	// User routes
	userHandler := user.NewHandler(user.Config{
		UserRepository:  s.userRepository,
		RegisterHandler: s.userRegisterHandler,
		LoginHandler:    s.userLoginQueryHandler,
	})
	userHandler.RegisterPublicRoutes(publicRouter)
	userHandler.RegisterProtectedRoutes(protectedRouter)

	// Expense routes
	expenseHandler := expense.NewHandler(s.addExpenseCommandHandler, s.timeService)
	expenseHandler.RegisterPublicRoutes(publicRouter)
	expenseHandler.RegisterProtectedRoutes(protectedRouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}

