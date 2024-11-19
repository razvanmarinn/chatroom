package postgres_repo

import (
	"github.com/google/uuid"
	_db "github.com/razvanmarinn/chatroom/internal/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	Db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		Db: db,
	}
}

func (pgr *PostgresUserRepository) GetUserByID(userUUID uuid.UUID) (_db.User, error) {
	var user _db.User
	result := pgr.Db.Where("id = ?", userUUID.String()).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return user, nil
		}
	}
	return user, nil
}

func (pgr *PostgresUserRepository) CreateUser(username string, password string) (uuid.UUID, error) {
	encryptedPassword, err := HashPassword(password)
	if err != nil {

		return uuid.Nil, err
	}

	newUser := _db.User{
		ID:       uuid.New(),
		Username: username,
		Password: encryptedPassword,
	}

	if err := pgr.Db.Create(&newUser).Error; err != nil {

		return uuid.Nil, err
	}

	return newUser.ID, nil
}

func (pgr *PostgresUserRepository) DeleteUser(username string) error {
	return nil
}

func (pgr *PostgresUserRepository) UserExists(username string) bool {
	var existingUser _db.User
	result := pgr.Db.Where("username = ?", username).First(&existingUser)

	if result.Error == nil {

		return true
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
