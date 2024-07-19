package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/beka-birhanu/finance-go/infrastructure"
)

func main() {
	// Load environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Initialize database connection
	db, err := infrastructure.NewPostgresDB(dbUser, dbPassword, dbName, dbHost, dbPort)
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
