package services

import (
	"fmt"

	"github.com/google/uuid"
	r_fact "github.com/razvanmarinn/chatroom/internal/db/repository_factory"
	"github.com/razvanmarinn/chatroom/internal/logger"

	cache "github.com/razvanmarinn/chatroom/internal/cache"
	"github.com/razvanmarinn/chatroom/internal/db"
)

type RoomService struct {
	RoomRepo     db.RoomRepository
	CacheManager cache.CacheManager
	Logger       logger.Logger
}

func NewRoomService(cm cache.CacheManager, rf r_fact.RepositoryFactory, logger logger.Logger) *RoomService {
	roomRepo, err := rf.CreateRoomRepository()
	if err != nil {
		fmt.Println("ASD")
	}
	return &RoomService{
		RoomRepo:     roomRepo,
		CacheManager: cm,
		Logger:       logger,
	}
}

func (rs *RoomService) RoomExists(roomName string) bool {
	return rs.RoomRepo.RoomExists(roomName)

}
func (rs *RoomService) GetRoomByName(roomName string) (db.Room, error) {
	return rs.RoomRepo.GetRoomByName(roomName)

}
func (rs *RoomService) CreateRoom(roomName string, owner_uuid uuid.UUID) (uuid.UUID, error) {
	return rs.RoomRepo.CreateRoom(roomName, owner_uuid)

}
