package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	instance *sql.DB   // Singleton instance of the database connection
	once     sync.Once // Ensures that the database connection is created only once
)

// Config holds the database connection configuration.
type Config struct {
	DbUser     string // Username for the database connection
	DbPassword string // Password for the database connection
	DbName     string // Name of the database
	DbHost     string // Host where the database server is located
	DbPort     string // Port on which the database server is listening
}

// Connect initializes and returns a singleton database connection.
func Connect(config Config) *sql.DB {
	once.Do(func() {
		var err error
		// Format the connection string using the provided configuration.
		connStr := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DbUser,
			config.DbPassword,
			config.DbHost,
			config.DbPort,
			config.DbName,
		)

		instance, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Could not connect to the database: %v", err)
		}

		err = instance.Ping()
		if err != nil {
			log.Fatalf("Could not ping the database: %v", err)
		}

		log.Println("DB: Successfully connected!")
	})

	return instance
}
