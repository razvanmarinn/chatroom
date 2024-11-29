package repo_factory

import (
	"fmt"

	"github.com/razvanmarinn/chatroom/internal/cfg"
	_db "github.com/razvanmarinn/chatroom/internal/db"

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
