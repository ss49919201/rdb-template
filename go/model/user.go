package model

import "time"

type User struct {
	ID        string `xorm:"id"`
	Name      string
	Count     int
	UpdatedAt time.Time
}
