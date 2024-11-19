package repo_factory

import (
	_db "github.com/razvanmarinn/chatroom/internal/db"
	pr "github.com/razvanmarinn/chatroom/internal/db/postgres"

	"gorm.io/gorm"
)

type PostgreSQLRepositoryFactory struct {
	db *gorm.DB
}

func NewPostgreSQLRepositoryFactory(db *gorm.DB) RepositoryFactory {
	return &PostgreSQLRepositoryFactory{db: db}
}
func (p *PostgreSQLRepositoryFactory) CreateUserRepository() (_db.UserRepository, error) {
	return pr.NewPostgresUserRepository(p.db), nil
}

func (p *PostgreSQLRepositoryFactory) CreateMessageRepository() (_db.MessageRepository, error) {
	return pr.NewPostgresMessageRepository(p.db), nil
}

func (p *PostgreSQLRepositoryFactory) CreateRoomRepository() (_db.RoomRepository, error) {
	return pr.NewPostgresRoomRepository(p.db), nil
}
