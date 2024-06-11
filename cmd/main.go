package main

import (
	"database/sql"
	"log"

	"github.com/gitKashish/ecommerce-api-go/cmd/api"
	"github.com/gitKashish/ecommerce-api-go/config"
	"github.com/gitKashish/ecommerce-api-go/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// Creating a new DB instance with configs from `config.Env`.
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAdress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Establishing Connecting with DB.
	initStorage(db)

	// Creating an Starting a new HTTP server.
	server := api.NewAPIServer(":8080", db)
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Ping - Check DB connection or establish one if not already established.
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB : Successfully connected.")
}
