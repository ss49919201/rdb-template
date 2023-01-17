package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	driverMySQL "github.com/go-sql-driver/mysql"
	"github.com/s-beats/rdb-template/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	_ = &model.User{}
	_ = newDB()
}

// parseTime=true が無いと Scan error on column index 3, name "updated_at": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
func newDB() *gorm.DB {
	// parseTime=true が無いと Scan error on column index 3, name "updated_at": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
	config := &driverMySQL.Config{
		User:      "user",
		Passwd:    "password",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "rdb",
		ParseTime: true,
	}
	db, err := gorm.Open(mysql.Open(config.FormatDSN()), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func insertRondom(ctx context.Context, db *gorm.DB) error {
	// 実行する毎に異なる動作が必要な場合は`rand.Seed(time.Now().UnixNano())`
	id := strconv.Itoa(rand.Intn(100))
	if err := db.WithContext(ctx).
		Table("users").
		Create(&model.User{
			ID:        id,
			Name:      fmt.Sprintf("name-%s", id),
			UpdatedAt: time.Unix(int64(rand.Intn(100)), 0),
		}).Error; err != nil {
		return err
	}
	return nil
}
