package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ss49919201/rdb-template/config"
)

func TestGetUserForUpdate(t *testing.T) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		t.Fatal(err)
	}

	getForUpdateWithWait := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUserForUpdate(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		time.Sleep(time.Second * 3)
		tx.Commit()
	}

	// expect error
	getForShare := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUserForShare(tx, "A"); err == nil {
			t.Error("expect error")
			return
		}
		tx.Commit()
	}

	get := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUser(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		tx.Commit()
	}

	go getForUpdateWithWait()
	time.Sleep(time.Second)
	get()
	getForShare()
}

func TestGetUserForShare(t *testing.T) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		t.Fatal(err)
	}

	getForShareWithWait := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUserForShare(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		time.Sleep(time.Second * 3)
		tx.Commit()
	}

	// expect error
	getForUpdate := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUserForUpdate(tx, "A"); err == nil {
			t.Error("expect error")
			return
		}
		tx.Commit()
	}

	getForShare := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUserForShare(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		tx.Commit()
	}

	get := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUser(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		tx.Commit()
	}

	go getForShareWithWait()
	time.Sleep(time.Second)
	get()
	getForUpdate()
	getForShare()
}

func TestGetUser(t *testing.T) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		t.Fatal(err)
	}

	getWithWait := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUser(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		time.Sleep(time.Second * 3)
		tx.Commit()
	}

	getForUpdate := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUserForUpdate(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		tx.Commit()
	}

	getForShare := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUserForShare(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		tx.Commit()
	}

	get := func() {
		tx, err := db.Begin()
		if err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		if _, err := getUser(tx, "A"); err != nil {
			t.Errorf("not expected: %s", err)
			return
		}
		tx.Commit()
	}

	go getWithWait()
	time.Sleep(time.Second)
	getForUpdate()
	getForShare()
	get()
}
