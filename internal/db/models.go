package db

import "time"

type User struct {
	Id        int       `gorm:"primaryKey"`
	Username  string    `gorm:"unique"`
	Password  string   
	CreatedAt time.Time `gorm:"autoCreateTime"` 
}