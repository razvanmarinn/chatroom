package db

import (
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByID(userID uuid.UUID) (User, error)
	UserExists(username string) bool
	CreateUser(username string, password string) (uuid.UUID, error)
	DeleteUser(username string) error
	GetUserByUsername(username string) (User, error)
}

type RoomRepository interface {
	GetRoomByName(roomName string) (Room, error)
	CreateRoom(roomName string, owner_uuid uuid.UUID) (uuid.UUID, error)
	RoomExists(roomName string) bool
}

type MessageRepository interface {
	CreateMessage(roomId uuid.UUID, userId uuid.UUID, content []byte) (*Message, error)
}
