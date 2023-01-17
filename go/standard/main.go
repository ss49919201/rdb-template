package main

import (
	"database/sql"
	"log"

	"github.com/ss49919201/rdb-template/config"
)

func main() {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := ping(db); err != nil {
		log.Fatal(err)
	}
}

func ping(db *sql.DB) error {
	return db.Ping()
}

// TODO
func getUser()
func getUserForUpdate()
