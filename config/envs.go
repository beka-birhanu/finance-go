package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost             string
	ServerPort             string
	DBHost                 string
	DBPort                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
	TestDBHost             string
	TestDBPort             string
	TestDBUser             string
	TestDBPassword         string
	TestDBName             string
}

var Envs = initConfig()

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

// Gets the env by key or fallbacks
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

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

