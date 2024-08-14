// Package config provides functionality to load and manage application configuration from environment variables.
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application's configuration values.
type Config struct {
	ServerHost             string // Hostname or IP address for the server
	ServerPort             string // Port number for the server
	DBHost                 string // Hostname or IP address for the database
	DBPort                 string // Port number for the database
	DBUser                 string // Username for the database
	DBPassword             string // Password for the database
	DBName                 string // Name of the database
	JWTSecret              string // Secret key for JWT signing
	JWTExpirationInSeconds int64  // JWT expiration time in seconds
	TestDBHost             string // Hostname or IP address for the test database
	TestDBPort             string // Port number for the test database
	TestDBUser             string // Username for the test database
	TestDBPassword         string // Password for the test database
	TestDBName             string // Name of the test database
}

// Envs holds the application's configuration loaded from environment variables.
var Envs = initConfig()

// initConfig initializes and returns the application configuration.
// It loads environment variables from a .env file and sets default values if needed.
func initConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Panicln(err)
	}

	return Config{
		ServerHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		ServerPort:             getEnv("PORT", "8080"),
		DBHost:                 getEnv("DB_HOST", "romareo"),
		DBPort:                 getEnv("DB_PORT", "5432"),
		DBUser:                 getEnv("DB_USER", "romareo"),
		DBPassword:             getEnv("DB_PASSWORD", "PythonIsTheGOAT"),
		DBName:                 getEnv("DB_NAME", "finance"),
		JWTSecret:              getEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 60*24),
		TestDBHost:             getEnv("TEST_DB_HOST", "localhost"),
		TestDBPort:             getEnv("TEST_DB_PORT", "5432"),
		TestDBUser:             getEnv("TEST_DB_USER", "test_user"),
		TestDBPassword:         getEnv("TEST_DB_PASSWORD", "test_password"),
		TestDBName:             getEnv("TEST_DB_NAME", "test_finance"),
	}
}

// getEnv retrieves the value of an environment variable or returns a fallback value if the variable is not set.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves the value of an environment variable as an integer or returns a fallback value if the variable is not set or cannot be parsed.
func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
