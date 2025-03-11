package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// InitDB инициализирует подключение к PostgreSQL
func InitDB() *sqlx.DB {
	dsn := "postgres://auth_user:auth_pass@localhost:5432/auth_db?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	log.Println("DB connection successful")
	return db
}
