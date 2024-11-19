package postgres_repo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/razvanmarinn/chatroom/internal/db"
	_db "github.com/razvanmarinn/chatroom/internal/db"
	"gorm.io/gorm"
)

type PostgresRoomRepository struct {
	Db *gorm.DB
}

func NewPostgresRoomRepository(db *gorm.DB) *PostgresRoomRepository {
	return &PostgresRoomRepository{
		Db: db,
	}
}
func (pgr *PostgresRoomRepository) CreateRoom(roomName string, owner_uuid uuid.UUID) (uuid.UUID, error) {
	room := &db.Room{
		ID:       uuid.New(),
		RoomName: roomName,
		Owner:    owner_uuid,
	}
	if err := pgr.Db.Create(&room).Error; err != nil {
		return uuid.Nil, err
	}
	return room.ID, nil

}

func (pgr *PostgresRoomRepository) RoomExists(roomName string) bool {
	var room _db.Room
	result := pgr.Db.Where("room_name = ?", roomName).First(&room)
	if result.Error != nil {

		return false
	}
	return true
}
func (pgr *PostgresRoomRepository) GetRoomByName(roomName string) (_db.Room, error) {
	var room _db.Room
	result := pgr.Db.Where("room_name = ?", roomName).First(&room)
	if result.Error != nil {
		fmt.Println("ASd") //TODO: Error improvements
	}
	return room, nil
}
