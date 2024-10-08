package db

import (
	"database/sql"
	"fmt"
	"log"

	"music_catalog/config"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// NewDB initializes a new database connection
func NewDB(config *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("[INFO] Successfully connected to the database")
	return db, nil
}
