package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var testDSN string

func TestMain(m *testing.M) {
	dsn, shutdown, err := setupDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	testDSN = dsn

	m.Run()

	fmt.Println("shutdown...")
	shutdown()
}

func setupDB() (dsn string, shutdown func(), err error) {
	// Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", nil, err
	}

	// Run
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "rdb-test",
		Repository: "mysql",
		Tag:        "5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=password",
			"MYSQL_DATABASE=rdb-test",
		},
		Cmd: []string{
			"mysqld", "--character-set-server=utf8mb4", "--collation-server=utf8mb4_unicode_ci",
		},
	})
	if err != nil {
		return "", nil, err
	}
	shutdown = func() {
		_ = resource.Close()
	}
	dsn = (&mysql.Config{
		User:                 "root",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 resource.GetHostPort("3306/tcp"),
		DBName:               "rdb-test",
		ParseTime:            true,
		AllowNativePasswords: true,
	}).FormatDSN()

	// リクエストを捌く準備ができるまで待つ
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", nil, err
	}
	if err := pool.Retry(func() error {
		return db.Ping()
	}); err != nil {
		shutdown()
		return "", nil, err
	}

	return
}
