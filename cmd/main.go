package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/beka-birhanu/finance-go/api/middleware"
	registercmd "github.com/beka-birhanu/finance-go/application/authentication/command"
	loginqry "github.com/beka-birhanu/finance-go/application/authentication/query"

	// expenseQuery "github.com/beka-birhanu/finance-go/application/expense/query"
	"github.com/beka-birhanu/finance-go/config"
	"github.com/beka-birhanu/finance-go/infrastructure/db"
	"github.com/beka-birhanu/finance-go/infrastructure/hash"
	"github.com/beka-birhanu/finance-go/infrastructure/jwt"
	"github.com/beka-birhanu/finance-go/infrastructure/repository"
	timeservice "github.com/beka-birhanu/finance-go/infrastructure/time_service"
)

var dbUser = config.Envs.DBUser
var dbPassword = config.Envs.DBPassword
var dbName = config.Envs.DBName
var dbHost = config.Envs.DBHost
var dbPort = config.Envs.DBPort

func main() {
	// Connect to the database
	database := db.Connect(db.Config{
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbName:     dbName,
		DbHost:     dbHost,
		DbPort:     dbPort,
	})

	// Initialize dependencies
	timeService := timeservice.New()

	userRepository := repository.NewUserRepository(database)

	jwtService := jwt.New(
		jwt.Config{
			SecretKey:   config.Envs.JWTSecret,
			Issuer:      config.Envs.ServerHost,
			ExpTime:     time.Duration(config.Envs.JWTExpirationInSeconds) * time.Second,
			TimeService: timeService,
		})
	hashService := hash.SingletonService()

	authorizationMiddleware := middleware.AuthorizationMiddleware(jwtService)

	// Initialize command and query handlers
	userRegisterCommandHandler := registercmd.NewHandler(registercmd.Config{
		UserRepository: userRepository,
		JwtService:     jwtService,
		HashService:    hashService,
		TimeService:    timeService,
	})

	userLoginQueryHandler := loginqry.NewHandler(loginqry.Config{
		UserRepository: userRepository,
		JwtService:     jwtService,
		HashService:    hashService,
	})

	// addExpenseHandler := expense.NewHandler(userRepository, timeService)

	// Create and run the server
	server := api.NewAPIServer(
		fmt.Sprintf(":%s", config.Envs.ServerPort),
		userRepository,
		userRegisterCommandHandler,
		userLoginQueryHandler,
		authorizationMiddleware,
		// addExpenseHandler,
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
