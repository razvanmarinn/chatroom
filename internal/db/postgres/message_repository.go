package postgres_repo

import (
	"github.com/google/uuid"
	_db "github.com/razvanmarinn/chatroom/internal/db"
	"gorm.io/gorm"
)

type PostgresMessageRepository struct {
	Db *gorm.DB
}

func NewPostgresMessageRepository(db *gorm.DB) *PostgresMessageRepository {
	return &PostgresMessageRepository{
		Db: db,
	}
}

func (pgr *PostgresMessageRepository) CreateMessage(roomId uuid.UUID, userId uuid.UUID, content []byte) (*_db.Message, error) {
	message := &_db.Message{
		RoomID:  roomId,
		UserID:  userId,
		Content: string(content),
	}

	return message, nil

}

func (pgr *PostgresMessageRepository) GetLastMessagesByRoomID(roomID uuid.UUID, numberOfMessages int) ([]_db.Message, error) {
	var messages []_db.Message
	err := pgr.Db.
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(numberOfMessages).
		Find(&messages).Error

	if err != nil {
		return messages, err
	}

	return messages, nil
}
