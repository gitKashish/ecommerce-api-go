package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// Load config and return a sql.DB instance with those configs.
func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
