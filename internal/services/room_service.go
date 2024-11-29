package services

import (
	"fmt"

	r_fact "github.com/razvanmarinn/chatroom/internal/db/repository_factory"

	cache "github.com/razvanmarinn/chatroom/internal/cache"
	"github.com/razvanmarinn/chatroom/internal/db"
)

type RoomService struct {
	RoomRepo     db.RoomRepository
	CacheManager cache.CacheManager
}

func NewRoomService(cm cache.CacheManager, rf r_fact.RepositoryFactory) *RoomService {
	roomRepo, err := rf.CreateRoomRepository()
	if err != nil {
		fmt.Println("ASD")
	}
	return &RoomService{
		RoomRepo:     roomRepo,
		CacheManager: cm,
	}
}

func (rs *RoomService) RoomExists(roomName string) bool {
	return rs.RoomRepo.RoomExists(roomName)

}
func (rs *RoomService) GetRoomByName(roomName string) (db.Room, error) {
	return rs.RoomRepo.GetRoomByName(roomName)

}
