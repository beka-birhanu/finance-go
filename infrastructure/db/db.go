package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	instance *sql.DB
	once     sync.Once
)

// Connect initializes and returns a singleton database connection
func Connect(dbUser, dbPassword, dbName, dbHost, dbPort string) *sql.DB {
	once.Do(func() {
		var err error
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
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
