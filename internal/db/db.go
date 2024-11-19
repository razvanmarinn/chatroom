package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("user=%s password=%s dbname=%s  host=localhost port=6543 sslmode=disable TimeZone=Europe/Paris", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Room{})

	Migrate(DB)

}

func Migrate(db *gorm.DB) error {
	err := db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id BIGSERIAL,
            room_id UUID,
            user_id UUID,
            content TEXT,
            created_at TIMESTAMP WITH TIME ZONE,
            PRIMARY KEY (id, created_at)
        ) PARTITION BY RANGE (created_at);
    `).Error

	if err != nil {
		return err
	}

	currentMonth := time.Now()
	for i := 0; i < 3; i++ {
		monthToCreate := currentMonth.AddDate(0, i, 0)
		err = CreateMonthlyPartition(db, monthToCreate.Format("2006-01"))
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateMonthlyPartition(db *gorm.DB, yearMonth string) error {
	return db.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS "messages_%s" 
        PARTITION OF messages 
        FOR VALUES FROM ('%s-01') TO ('%s-01')
        WITH (fillfactor = 70);
    `, yearMonth, yearMonth, nextMonth(yearMonth))).Error
}

func nextMonth(yearMonth string) string {
	t, err := time.Parse("2006-01", yearMonth)
	if err != nil {
		return time.Now().Format("2006-01")
	}

	nextMonthTime := t.AddDate(0, 1, 0)

	return nextMonthTime.Format("2006-01")
}

func CheckRoomExists(roomName string) bool {
	var room Room
	result := DB.Where("room_name = ?", roomName).First(&room)
	if result.Error != nil {

		return false
	}
	return true
}
func GetRoomUUID(roomName string) (uuid.UUID, error) {
	var room Room
	result := DB.Where("room_name = ?", roomName).First(&room)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return uuid.Nil, nil
		}
		return uuid.Nil, result.Error
	}
	return room.ID, nil
}

func GetUsername(userUUID uuid.UUID) (string, error) {
	var user User
	result := DB.Where("id = ?", userUUID.String()).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", nil
		}
	}
	return user.Username, nil
}
