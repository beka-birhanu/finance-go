package api

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/expenses"
	"github.com/beka-birhanu/finance-go/api/users"
	commandAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/authentication"
	"github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/expense"
	querieAuth "github.com/beka-birhanu/finance-go/application/common/cqrs/i_queries/authentication"
	"github.com/beka-birhanu/finance-go/application/common/interfaces/persistance"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr                       string
	userRepository             persistance.IUserRepository
	userRegisterCommandHandler commandAuth.IUserRegisterCommandHandler
	userLoginQueryHandler      querieAuth.IUserLoginQueryHandler
	authorizationMiddleware    func(http.Handler) http.Handler
	addExpenseCommandHandler   expense.IAddExpenseCommand
}

func NewAPIServer(
	addr string,
	userRepository persistance.IUserRepository,
	userRegisterCommandHandler commandAuth.IUserRegisterCommandHandler,
	userQueryHandler querieAuth.IUserLoginQueryHandler,
	authorizationMiddleware func(http.Handler) http.Handler,
	addExpenseCommandHandler expense.IAddExpenseCommand,
) *APIServer {
	return &APIServer{
		addr:                       addr,
		userRepository:             userRepository,
		userRegisterCommandHandler: userRegisterCommandHandler,
		userLoginQueryHandler:      userQueryHandler,
		authorizationMiddleware:    authorizationMiddleware,
		addExpenseCommandHandler:   addExpenseCommandHandler,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// Routes that do not require authentication
	publicRouter := router.PathPrefix("/api/v1/public").Subrouter()

	// Routes that require authentication
	protectedRouter := router.PathPrefix("/api/v1").Subrouter()
	protectedRouter.Use(s.authorizationMiddleware)

	userHandler := users.NewHandler(
		s.userRepository,
		s.userRegisterCommandHandler,
		s.userLoginQueryHandler,
	)
	userHandler.RegisterPublicRoutes(publicRouter)
	userHandler.RegisterProtectedRoutes(protectedRouter)

	expenseHandler := expenses.NewHandler(s.addExpenseCommandHandler)
	expenseHandler.RegisterPublicRoutes(publicRouter)
	expenseHandler.RegisterProtectedRoutes(protectedRouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
