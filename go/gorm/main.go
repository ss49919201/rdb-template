package main

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := gorm.Open(mysql.Open("user:password@tcp(localhost:3306)/rdb"), &gorm.Config{
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

	ctx := context.Background()
	insertWithTx(ctx, db, "1")
}

func insertWithTx(ctx context.Context, db *gorm.DB, key string) {
	type User struct {
		ID        string `xorm:"id"`
		Name      string
		Count     int
		UpdatedAt time.Time
	}

	tx := db.WithContext(ctx)
	if err := tx.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&User{ID: key}).
			Count(&count).
			Error; err != nil {
			return err
		}

		if count > 0 {
			panic("already exsit")
		}

		return tx.Create(&User{
			ID:        key,
			Name:      "name",
			Count:     1,
			UpdatedAt: time.Now(),
		}).Error
	}); err != nil {
		panic(err)
	}
}

func newMySQLClient() (*gorm.DB, error) {
	return gorm.Open(mysql.Open("user:password@tcp(localhost:3306)/rdb)"), &gorm.Config{
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
}
