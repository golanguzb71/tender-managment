package db

import (
	"database/sql"
	"fmt"
	"log"
	"tender-managment/internal/config"

	_ "github.com/lib/pq"
)

func NewDatabase(databaseConf *config.DatabaseConfig) *sql.DB {

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		databaseConf.Username,
		databaseConf.Password,
		databaseConf.Host,
		databaseConf.Port,
		databaseConf.DBName,
		"disable",
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
