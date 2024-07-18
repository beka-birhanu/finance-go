package infrastructure

// import "github.com/beka-birhanu/finance-go/application"

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresDB(user, password, dbname, host, port string) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	return sql.Open("postgres", connStr)
}
