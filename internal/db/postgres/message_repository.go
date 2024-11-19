package postgres_repo

import (
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

func (pgr *PostgresMessageRepository) CreateMessage(username string, password string) (int64, error) {
	return 10, nil
}
