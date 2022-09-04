package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/s-beats/rdb-template/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// parseTime=true が無いと Scan error on column index 3, name "updated_at": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
	db, err := gorm.Open(mysql.Open("user:password@tcp(localhost:3306)/rdb?parseTime=true"), &gorm.Config{
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
	decreaseCount(ctx, db, "1")
}

func getOne(ctx context.Context, db *gorm.DB, key string) (*model.User, error) {
	var u *model.User
	if err := db.WithContext(ctx).
		Table("users").
		Where("id = ?", key).
		Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func exist(ctx context.Context, db *gorm.DB, key string) (bool, error) {
	var n int64
	if err := db.WithContext(ctx).
		Table("users").
		Where("id = ?", key).
		Count(&n).Error; err != nil {
		return false, err
	}

	return n > 0, nil
}

func insertWithTx(ctx context.Context, db *gorm.DB, key string) {
	tx := db.WithContext(ctx)
	if err := tx.Transaction(func(tx *gorm.DB) error {
		exist, err := exist(ctx, tx, key)
		if err != nil {
			return err
		}

		if exist {
			panic("already exsit")
		}

		return tx.Create(&model.User{
			ID:        key,
			Name:      "name",
			Count:     1,
			UpdatedAt: time.Now(),
		}).Error
	}); err != nil {
		panic(err)
	}
}

func decreaseCount(ctx context.Context, db *gorm.DB, key string) {
	tx := db.WithContext(ctx)
	if err := tx.Transaction(func(tx *gorm.DB) error {
		var u *model.User
		if err := tx.Where("id = ?", key).
			Find(&u).Error; err != nil {
			return err
		}

		if u == nil {
			panic("not exists")
		}

		return tx.Updates(&model.User{
			ID:        key,
			Name:      "name",
			Count:     u.Count - 1,
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
