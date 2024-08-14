package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/beka-birhanu/finance-go/api/expense"
	"github.com/beka-birhanu/finance-go/api/middleware"
	"github.com/beka-birhanu/finance-go/api/router"
	"github.com/beka-birhanu/finance-go/api/user"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	"github.com/beka-birhanu/finance-go/config"
	"github.com/beka-birhanu/finance-go/infrastructure/db"
	"github.com/beka-birhanu/finance-go/infrastructure/hash"
	"github.com/beka-birhanu/finance-go/infrastructure/jwt"
	expenserepo "github.com/beka-birhanu/finance-go/infrastructure/repository/expense"
	userrepo "github.com/beka-birhanu/finance-go/infrastructure/repository/user"
	timeservice "github.com/beka-birhanu/finance-go/infrastructure/time_service"
)

// Global variables to hold database configuration
var (
	dbUser     = config.Envs.DBUser
	dbPassword = config.Envs.DBPassword
	dbName     = config.Envs.DBName
	dbHost     = config.Envs.DBHost
	dbPort     = config.Envs.DBPort
)

func main() {
	// Connect to the database
	database := db.Connect(db.Config{
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbName:     dbName,
		DbHost:     dbHost,
		DbPort:     dbPort,
	})

	// Initialize services and repositories
	timeService := timeservice.New()
	userRepository := userrepo.New(database)
	expenseRepository := expenserepo.New(database)
	jwtService := initializeJWTService(timeService)
	hashService := hash.SingletonService()
	authorizationMiddleware := middleware.Authorization(jwtService)

	// Initialize command and query handlers
	userRegisterCommandHandler := initializeUserRegisterHandler(userRepository, jwtService, hashService, timeService)
	userLoginQueryHandler := initializeUserLoginQueryHandler(userRepository, jwtService, hashService)
	addExpenseHandler := initializeAddExpenseHandler(userRepository, timeService)
	getExpenseHandler := initializeGetExpenseHandler(expenseRepository)
	getExpensesHandler := initializeGetExpensesHandler(expenseRepository)
	patchExpenseHandler := initializePatchExpenseHandler(expenseRepository)

	userHandler := user.NewHandler(user.Config{
		UserRepository:  userRepository,
		RegisterHandler: userRegisterCommandHandler,
		LoginHandler:    userLoginQueryHandler,
	})

	// Expense routes
	expenseHandler := expense.NewHandler(expense.Config{
		AddHandler:         addExpenseHandler,
		GetHandler:         getExpenseHandler,
		PatchHandler:       patchExpenseHandler,
		GetMultipleHandler: getExpensesHandler,
	})

	// Create and run the server
	server := router.NewRouter(router.Config{
		Addr:                    fmt.Sprintf(":%s", config.Envs.ServerPort),
		Controllers:             []api.IController{userHandler, expenseHandler},
		AuthorizationMiddleware: authorizationMiddleware,
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initializePatchExpenseHandler(expenseRepository *expenserepo.Repository) *expensecmd.PatchHandler {
	return expensecmd.NewPatchHandler(expenseRepository)
}

func initializeGetExpenseHandler(expenseRepository *expenserepo.Repository) *expensqry.GetHandler {
	return expensqry.NewGetHandler(expenseRepository)
}

func initializeGetExpensesHandler(expenseRepository *expenserepo.Repository) *expensqry.GetMultipleHandler {
	return expensqry.NewGetMultipleHandler(expenseRepository)
}

// initializeJWTService initializes and returns a new JWT service instance.
func initializeJWTService(timeService *timeservice.Service) *jwt.Service {
	return jwt.New(
		jwt.Config{
			SecretKey:   config.Envs.JWTSecret,
			Issuer:      config.Envs.ServerHost,
			ExpTime:     time.Duration(config.Envs.JWTExpirationInSeconds) * time.Second,
			TimeService: timeService,
		})
}

// initializeUserRegisterHandler initializes and returns a new user register command handler.
func initializeUserRegisterHandler(userRepo *userrepo.Repository, jwtService *jwt.Service, hashService *hash.Service, timeService *timeservice.Service) *registercmd.Handler {
	return registercmd.NewHandler(registercmd.Config{
		UserRepo: userRepo,
		JwtSvc:   jwtService,
		HashSvc:  hashService,
		TimeSvc:  timeService,
	})
}

// initializeUserLoginQueryHandler initializes and returns a new user login query handler.
func initializeUserLoginQueryHandler(userRepo *userrepo.Repository, jwtService *jwt.Service, hashService *hash.Service) *loginqry.Handler {
	return loginqry.NewHandler(loginqry.Config{
		UserRepository: userRepo,
		JwtService:     jwtService,
		HashService:    hashService,
	})
}

// initializeAddExpenseHandler initializes and returns a new add expense command handler.
func initializeAddExpenseHandler(userRepo *userrepo.Repository, timeService *timeservice.Service) *expensecmd.AddHandler {
	return expensecmd.NewAddHandler(expensecmd.Config{
		UserRepository: userRepo,
		TimeService:    timeService,
	})
}
