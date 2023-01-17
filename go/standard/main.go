package main

import (
	"database/sql"
	"log"

	"github.com/s-beats/rdb-template/config"
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
