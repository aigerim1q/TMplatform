package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	connStr := "host=localhost port=5432 user=silence_ai_backender password=FreshFreckles1987 dbname=stages_projects sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
