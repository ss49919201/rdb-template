package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/samber/lo"

	"github.com/ss49919201/rdb-template/config"
	"github.com/ss49919201/rdb-template/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db, err := newMySQLClient()
	if err != nil {
		panic("failed to connect database")
	}

	updateWithoutDefaultTX(context.Background(), db)
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

func updateWithoutDefaultTX(ctx context.Context, db *gorm.DB) {
	tx := db.Session(
		&gorm.Session{
			Context:                ctx,
			SkipDefaultTransaction: true,
		},
	)

	var user *model.User
	err := tx.First(&user, 1).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(err)
		}
	}
	if !lo.IsEmpty(*user) {
		if err := tx.Model(&user).Update("name", "HOGE").Error; err != nil {
			panic(err)
		}
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
			return errors.New("not found")
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
	return gorm.Open(mysql.Open(config.DSN()), &gorm.Config{
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
