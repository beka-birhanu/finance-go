package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/beka-birhanu/finance-go/application/authentication/commands"
	"github.com/beka-birhanu/finance-go/application/authentication/queries"
	"github.com/beka-birhanu/finance-go/configs"
	"github.com/beka-birhanu/finance-go/infrastructure/db"
	"github.com/beka-birhanu/finance-go/infrastructure/hash"
	"github.com/beka-birhanu/finance-go/infrastructure/jwt"
	"github.com/beka-birhanu/finance-go/infrastructure/repositories"
)

func main() {
	// Connect to the database
	database := db.Connect()

	// Initialize dependencies
	userRepository := repositories.NewUserRepository(database)
	jwtService := jwt.NewJwtService(
		configs.Envs.JWTSecret,
		configs.Envs.ServerHost,
		time.Duration(configs.Envs.JWTExpirationInSeconds)*time.Second,
	)
	hashService := hash.GetHashService()

	// Initialize command and query handlers
	userRegisterCommandHandler := commands.NewRegisterCommandHandler(userRepository, jwtService, hashService)
	userLoginQueryHandler := queries.NewUserLoginQueryHandler(userRepository, jwtService, hashService)

	// Create and run the server
	server := api.NewAPIServer(
		fmt.Sprintf(":%s", configs.Envs.ServerPort),
		userRepository,
		userRegisterCommandHandler,
		userLoginQueryHandler,
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
