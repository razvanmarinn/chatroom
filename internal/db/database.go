package db

import (
	"fmt"
	"log"
	"os"

	"github.com/razvanmarinn/chatroom/internal/cfg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfg cfg.Config) (*gorm.DB, error) {


	//TODO: Multiple databases
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=localhost port=6543 sslmode=disable TimeZone=Europe/Paris",
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
	)
	DB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	return DB, nil
}
