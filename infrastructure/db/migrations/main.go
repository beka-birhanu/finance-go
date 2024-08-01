package main

import (
	"fmt"
	"log"

	"github.com/beka-birhanu/finance-go/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	dbUser     = config.Envs.DBUser
	dbPassword = config.Envs.DBPassword
	dbName     = config.Envs.DBName
	dbHost     = config.Envs.DBHost
	dbPort     = config.Envs.DBPort
)

func main() {
	// Construct the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize the migrate instance
	m, err := migrate.New(
		"file://infrastructure/db/migrations", // Path to migration files
		connStr,                               // Connection string
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate instance: %v", err)
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Hey there, migration has failed: %v", err)
	}

	log.Println("Migration completed successfully")
}
