package main

import (
	"fmt"
	"log"

	"github.com/beka-birhanu/finance-go/api"
	"github.com/beka-birhanu/finance-go/configs"
	"github.com/beka-birhanu/finance-go/infrastructure/db"
)

func main() {
	db := db.Connect()

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.ServerPort), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
