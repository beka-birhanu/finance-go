package api

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/expense"
	"github.com/beka-birhanu/finance-go/api/user"
	authCommand "github.com/beka-birhanu/finance-go/application/authentication/command"
	"github.com/beka-birhanu/finance-go/application/authentication/common"
	authQuery "github.com/beka-birhanu/finance-go/application/authentication/query"
	handlerInterface "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	"github.com/beka-birhanu/finance-go/application/common/interface/repository"
	expenseCommand "github.com/beka-birhanu/finance-go/application/expense/command"
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr                       string
	userRepository             repository.IUserRepository
	userRegisterCommandHandler handlerInterface.ICommandHandler[*authCommand.UserRegisterCommand, *common.AuthResult]
	userLoginQueryHandler      handlerInterface.ICommandHandler[*authQuery.UserLoginQuery, *common.AuthResult]
	authorizationMiddleware    func(http.Handler) http.Handler
	addExpenseCommandHandler   handlerInterface.ICommandHandler[*expenseCommand.AddExpenseCommand, *model.Expense]
}

func NewAPIServer(
	addr string,
	userRepository repository.IUserRepository,
	userRegisterCommandHandler handlerInterface.ICommandHandler[*authCommand.UserRegisterCommand, *common.AuthResult],
	userLoginQueryHandler handlerInterface.ICommandHandler[*authQuery.UserLoginQuery, *common.AuthResult],
	authorizationMiddleware func(http.Handler) http.Handler,
	addExpenseCommandHandler handlerInterface.ICommandHandler[*expenseCommand.AddExpenseCommand, *model.Expense],
) *APIServer {
	return &APIServer{
		addr:                       addr,
		userRepository:             userRepository,
		userRegisterCommandHandler: userRegisterCommandHandler,
		userLoginQueryHandler:      userLoginQueryHandler,
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

	userHandler := user.NewHandler(
		s.userRepository,
		s.userRegisterCommandHandler,
		s.userLoginQueryHandler,
	)
	userHandler.RegisterPublicRoutes(publicRouter)
	userHandler.RegisterProtectedRoutes(protectedRouter)

	expenseHandler := expense.NewHandler(s.addExpenseCommandHandler)
	expenseHandler.RegisterPublicRoutes(publicRouter)
	expenseHandler.RegisterProtectedRoutes(protectedRouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
