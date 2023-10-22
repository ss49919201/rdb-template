package config

import "github.com/go-sql-driver/mysql"

// parseTime=true が無いと Scan error on column index 3, name "updated_at": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
func DSN() string {
	return (&mysql.Config{
		User:                 "user",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "rdb",
		ParseTime:            true,
		AllowNativePasswords: true,
	}).FormatDSN()
}
