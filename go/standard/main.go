package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/ss49919201/rdb-template/config"
	"github.com/ss49919201/rdb-template/model"
)

func main() {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := ping(db); err != nil {
		log.Fatal(err)
	}
	test(db)
}

func test(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	user, err := getUser(tx, "A")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println([]byte(user.Name))
	}
	time.Sleep(time.Second * 100)
	tx.Commit()
}

func ping(db *sql.DB) error {
	return db.Ping()
}

func getUser(tx *sql.Tx, id string) (*model.User, error) {
	user := &model.User{}
	var updatedAt sql.NullTime
	err := tx.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Count, &updatedAt)
	if err != nil {
		return nil, err
	}
	user.UpdatedAt = updatedAt.Time
	return user, nil
}

func getUserForShare(tx *sql.Tx, id string) (*model.User, error) {
	user := &model.User{}
	var updatedAt sql.NullTime
	err := tx.QueryRow("SELECT * FROM users WHERE id = ? FOR SHARE NOWAIT", id).Scan(&user.ID, &user.Name, &user.Count, &updatedAt)
	if err != nil {
		return nil, err
	}
	user.UpdatedAt = updatedAt.Time
	return user, nil
}

func getUserForUpdate(tx *sql.Tx, id string) (*model.User, error) {
	user := &model.User{}
	var updatedAt sql.NullTime
	err := tx.QueryRow("SELECT * FROM users WHERE id = ? FOR UPDATE NOWAIT", id).Scan(&user.ID, &user.Name, &user.Count, &updatedAt)
	if err != nil {
		return nil, err
	}
	user.UpdatedAt = updatedAt.Time
	return user, nil
}
