package config

import "github.com/go-sql-driver/mysql"

func DSN() string {
	return (&mysql.Config{
		User:      "user",
		Passwd:    "password",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "rdb",
		ParseTime: true,
	}).FormatDSN()
}
