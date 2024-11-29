package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	r_fact "github.com/razvanmarinn/chatroom/internal/db/repository_factory"

	cache "github.com/razvanmarinn/chatroom/internal/cache"
	"github.com/razvanmarinn/chatroom/internal/db"
)

type MessageService struct {
	MessageRepo  db.MessageRepository
	CacheManager cache.CacheManager
}

func NewMessageService(cm cache.CacheManager, rf r_fact.RepositoryFactory) *MessageService {
	messageRepo, err := rf.CreateMessageRepository()
	if err != nil {
		fmt.Println("ASD")
	}
	return &MessageService{
		MessageRepo:  messageRepo,
		CacheManager: cm,
	}
}

func (ms *MessageService) GetLastMessagesByRoomID(roomUUID uuid.UUID, numberOfMessages int) ([]db.Message, error) {
	ctx := context.Background()

	listKey := roomUUID.String() + "_messages"
	messageList, err := ms.CacheManager.GetList(ctx, listKey)
	if err != nil {
		return nil, err
	}

	// return messageList, nil
	return ms.MessageRepo.GetLastMessagesByRoomID(roomUUID, numberOfMessages)
}

func (ms *MessageService) CreateMessage(roomUUID uuid.UUID, userID uuid.UUID, content []byte) (*db.Message, error) {
	ctx := context.Background()
	listKey := roomUUID.String() + "_messages"
	message, err := ms.MessageRepo.CreateMessage(roomUUID, userID, content)
	length, err := ms.CacheManager.GetLengthForList(ctx, listKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get list length: %w", err)
	}
	if length > 100 {
		if err := ms.CacheManager.DeleteFromList(ctx, listKey); err != nil {
			return nil, fmt.Errorf("failed to delete list: %w", err)
		}
	}

	if err := ms.CacheManager.AddToList(ctx, listKey, string(content)); err != nil {
		return nil, fmt.Errorf("failed to add message to list: %w", err)
	}

	return message, err
}
