package main

import (
	"fmt"
	"log"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/beka-birhanu/finance-go/configs"
	"github.com/beka-birhanu/finance-go/infrastructure/db"
	"github.com/beka-birhanu/finance-go/infrastructure/repositories"
)

func main() {
	db := db.Connect()

	userRepository := repositories.NewUserRepository(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.ServerPort), userRepository)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

