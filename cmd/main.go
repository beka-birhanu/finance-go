package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/beka-birhanu/finance-go/api/middleware"
	authCommand "github.com/beka-birhanu/finance-go/application/authentication/command"
	authQuery "github.com/beka-birhanu/finance-go/application/authentication/query"
	expenseCommand "github.com/beka-birhanu/finance-go/application/expense/command"

	// expenseQuery "github.com/beka-birhanu/finance-go/application/expense/query"
	"github.com/beka-birhanu/finance-go/config"
	"github.com/beka-birhanu/finance-go/infrastructure/db"
	"github.com/beka-birhanu/finance-go/infrastructure/hash"
	"github.com/beka-birhanu/finance-go/infrastructure/jwt"
	"github.com/beka-birhanu/finance-go/infrastructure/repository"
)

var dbUser = config.Envs.DBUser
var dbPassword = config.Envs.DBPassword
var dbName = config.Envs.DBName
var dbHost = config.Envs.DBHost
var dbPort = config.Envs.DBPort

func main() {
	// Connect to the database
	database := db.Connect(dbUser, dbPassword, dbName, dbHost, dbPort)

	// Initialize dependencies
	userRepository := repository.NewUserRepository(database)
	jwtService := jwt.NewJwtService(
		config.Envs.JWTSecret,
		config.Envs.ServerHost,
		time.Duration(config.Envs.JWTExpirationInSeconds)*time.Second,
	)
	hashService := hash.GetHashService()
	authorizationMiddleware := middleware.AuthorizationMiddleware(jwtService)

	// Initialize command and query handlers
	userRegisterCommandHandler := authCommand.NewRegisterCommandHandler(userRepository, jwtService, hashService)
	userLoginQueryHandler := authQuery.NewUserLoginQueryHandler(userRepository, jwtService, hashService)
	addExpenseHandler := expenseCommand.NewAddExpenseCommandHandler(userRepository)

	// Create and run the server
	server := api.NewAPIServer(
		fmt.Sprintf(":%s", config.Envs.ServerPort),
		userRepository,
		userRegisterCommandHandler,
		userLoginQueryHandler,
		authorizationMiddleware,
		addExpenseHandler,
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
