package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/beka-birhanu/finance-go/api/middleware"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
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
	addExpenseHandler := initializeAddExpenseHandler(userRepository, timeService, expenseRepository)

	// Create and run the server
	server := api.NewAPIServer(api.Config{
		Addr:                     fmt.Sprintf(":%s", config.Envs.ServerPort),
		UserRepository:           userRepository,
		UserRegisterHandler:      userRegisterCommandHandler,
		UserLoginQueryHandler:    userLoginQueryHandler,
		AuthorizationMiddleware:  authorizationMiddleware,
		AddExpenseCommandHandler: addExpenseHandler,
		TimeService:              timeService,
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
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
		UserRepository: userRepo,
		JwtService:     jwtService,
		HashService:    hashService,
		TimeService:    timeService,
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
func initializeAddExpenseHandler(userRepo *userrepo.Repository, timeService *timeservice.Service, expenseRepo *expenserepo.Repository) *expensecmd.AddHandler {
	return expensecmd.NewAddHandler(expensecmd.Config{
		UserRepository:    userRepo,
		TimeService:       timeService,
		ExpenseRepository: expenseRepo,
	})
}

