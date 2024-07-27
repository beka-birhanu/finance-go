package api

import (
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/user"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	auth "github.com/beka-birhanu/finance-go/application/authentication/common"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	irepository "github.com/beka-birhanu/finance-go/application/common/interface/repository"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr                    string
	userRepository          irepository.IUserRepository
	userRegisterHandler     icmd.IHandler[*registercmd.Command, *auth.Result]
	userLoginQueryHandler   iquery.IHandler[*loginqry.Query, *auth.Result]
	authorizationMiddleware func(http.Handler) http.Handler
	// addExpenseCommandHandler   handlerInterface.ICommandHandler[*expenseCommand.AddExpenseCommand, *model.Expense]
}

func NewAPIServer(
	addr string,
	userRepository irepository.IUserRepository,
	userRegisterHandler icmd.IHandler[*registercmd.Command, *auth.Result],
	userLoginQueryHandler iquery.IHandler[*loginqry.Query, *auth.Result],
	authorizationMiddleware func(http.Handler) http.Handler,
	// addExpenseCommandHandler handlerInterface.ICommandHandler[*expenseCommand.AddExpenseCommand, *model.Expense],
) *APIServer {
	return &APIServer{
		addr:                    addr,
		userRepository:          userRepository,
		userRegisterHandler:     userRegisterHandler,
		userLoginQueryHandler:   userLoginQueryHandler,
		authorizationMiddleware: authorizationMiddleware,
		// addExpenseCommandHandler:   addExpenseCommandHandler,
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
		user.Config{
			UserRepository:  s.userRepository,
			RegisterHandler: s.userRegisterHandler,
			LoginHandler:    s.userLoginQueryHandler,
		})
	userHandler.RegisterPublicRoutes(publicRouter)
	userHandler.RegisterProtectedRoutes(protectedRouter)

	// expenseHandler := expense.NewHandler(s.addExpenseCommandHandler)
	// expenseHandler.RegisterPublicRoutes(publicRouter)
	// expenseHandler.RegisterProtectedRoutes(protectedRouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
