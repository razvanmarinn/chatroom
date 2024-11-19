package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID    `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Message struct {
    ID        uint64    `gorm:"primaryKey;column:id"`
    RoomID    uuid.UUID `gorm:"column:room_id"`
    UserID    uuid.UUID `gorm:"column:user_id"`
    Content   string    `gorm:"column:content"`
    CreatedAt time.Time `gorm:"column:created_at"`
}

type Room struct {
	ID uuid.UUID `gorm:"primaryKey;column:id"`
	RoomName string 
  	Owner     uuid.UUID `gorm:"column:owner_id;foreignKey:ID;references:users"`
	CreatedAt time.Time `gorm:"column:created_at"`
}