package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// InitDB инициализирует подключение к PostgreSQL
func InitDB() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	log.Println(dsn)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	log.Println("DB connection successful")
	return db
}
