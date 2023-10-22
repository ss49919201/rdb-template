package main

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestQueryWithPrepare(t *testing.T) {
	db, err := sqlx.Open("mysql", testDSN)
	if err != nil {
		assert.Fail(t, "failed to open database handle")
	}

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS tbl (id INT NOT NULL, PRIMARY KEY (id))"); err != nil {
		assert.Fail(t, "failed to create table")
	}
	if _, err := db.Exec("INSERT INTO tbl (id) VALUES(1)"); err != nil {
		assert.Fail(t, "failed to insert row")
	}
	defer func() {
		db.Exec("DROP TABLE IF EXISTS users")
	}()

	t.Run("Type is int64", func(t *testing.T) {
		result := QueryWithPrepare(db)
		if assert.Greater(t, len(result), 0) {
			assert.Equal(t, "int64", fmt.Sprintf("%T", result[0]["id"]))
		}
	})
}
