package main

import (
	"log"
	"os"

	"github.com/gitKashish/ecommerce-api-go/config"
	"github.com/gitKashish/ecommerce-api-go/db"
	mySqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mySqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// creating MySQL configuration object.
	// loading values from the environment.
	cfg := mySqlDriver.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAdress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// initiating a new DB Instance (Handle).
	db, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// creating a Migrate driver using the DB Instance.
	driver, err := mySqlMigrate.WithInstance(db, &mySqlMigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// creating migration object with ...
	// migrations path, driver name, & migration driver.
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// debug information.
	v, d, _ := m.Version()
	log.Printf("Version: %d, dirty: %v", v, d)

	// executing Up-migration or Down-migration as per CLI argument.
	// Command executed through `Makefile` in this case.
	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
