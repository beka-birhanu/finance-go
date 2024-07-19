package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/configs"
	"github.com/beka-birhanu/finance-go/infrastructure"
)

func main() {
	// Load environment variables
	dbUser := configs.Envs.DBUser
	dbPassword := configs.Envs.DBPassword
	dbName := configs.Envs.DBName
	dbHost := configs.Envs.DBHost
	dbPort := configs.Envs.DBPort

	// Initialize database connection
	db, err := infrastructure.NewDB(dbUser, dbPassword, dbName, dbHost, dbPort)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
