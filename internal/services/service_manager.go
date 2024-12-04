package services

import (
	cache "github.com/razvanmarinn/chatroom/internal/cache"
	r_fact "github.com/razvanmarinn/chatroom/internal/db/repository_factory"
	"github.com/razvanmarinn/chatroom/internal/logger"
)

type ServiceManager struct {
	UserService    *UserService
	RoomService    *RoomService
	MessageService *MessageService
}

func NewServiceManager(cacheManager cache.CacheManager, repoFactory r_fact.RepositoryFactory, logger logger.Logger) *ServiceManager {
	messageService := NewMessageService(cacheManager, repoFactory, logger)
	userService, _ := NewUserService(cacheManager, repoFactory, logger)
	roomService := NewRoomService(cacheManager, repoFactory, logger)

	serviceManager := &ServiceManager{
		UserService:    userService,
		MessageService: messageService,
		RoomService: roomService,
	}
	return serviceManager
}
