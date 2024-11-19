package repo_factory

import (
	"fmt"
	_db "github.com/razvanmarinn/chatroom/internal/db"
	"github.com/razvanmarinn/chatroom/internal/cfg"

	"gorm.io/gorm"
)



type RepositoryFactory interface {
	CreateUserRepository() (_db.UserRepository, error)
	CreateRoomRepository() (_db.RoomRepository, error)
	CreateMessageRepository() (_db.MessageRepository, error)
}

func CreateRepositoryFactory(dbType cfg.DatabaseType, db *gorm.DB) (RepositoryFactory, error) {
	switch dbType {
	case cfg.PostgreSQL:
		return NewPostgreSQLRepositoryFactory(db), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
