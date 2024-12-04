package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	r_fact "github.com/razvanmarinn/chatroom/internal/db/repository_factory"
	"github.com/razvanmarinn/chatroom/internal/logger"

	cache "github.com/razvanmarinn/chatroom/internal/cache"
	"github.com/razvanmarinn/chatroom/internal/db"
)

type MessageService struct {
	MessageRepo  db.MessageRepository
	CacheManager cache.CacheManager
	Logger       logger.Logger
}

func NewMessageService(cm cache.CacheManager, rf r_fact.RepositoryFactory, logger logger.Logger) *MessageService {
	messageRepo, err := rf.CreateMessageRepository()
	if err != nil {
		fmt.Println("ASD")
	}
	return &MessageService{
		MessageRepo:  messageRepo,
		CacheManager: cm,
		Logger:       logger,
	}
}

func (ms *MessageService) GetLastMessagesByRoomID(roomUUID uuid.UUID, numberOfMessages int) ([]db.Message, error) {
	ctx := context.Background()

	listKey := roomUUID.String() + "_messages"
	messageListStr, err := ms.CacheManager.GetList(ctx, listKey)
	if err != nil {
		return nil, err
	}
	messageList := []db.Message{}

	for _, messageStr := range messageListStr {
		var message db.Message
		if err := json.Unmarshal([]byte(messageStr), &message); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
		messageList = append(messageList, message)
	}

	if len(messageList) > 0 {
		return messageList, nil
	}

	return ms.MessageRepo.GetLastMessagesByRoomID(roomUUID, numberOfMessages)
}

func (ms *MessageService) CreateMessage(roomUUID uuid.UUID, userID uuid.UUID, content []byte) (*db.Message, error) {
	ctx := context.Background()
	listKey := roomUUID.String() + "_messages"
	message, err := ms.MessageRepo.CreateMessage(roomUUID, userID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message to JSON: %w", err)
	}
	length, err := ms.CacheManager.GetLengthForList(ctx, listKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get list length: %w", err)
	}
	if length > 100 {
		if err := ms.CacheManager.DeleteFromList(ctx, listKey); err != nil {
			return nil, fmt.Errorf("failed to delete list: %w", err)
		}
	}

	if err := ms.CacheManager.AddToList(ctx, listKey, string(messageJSON)); err != nil {
		return nil, fmt.Errorf("failed to add message to list: %w", err)
	}

	return message, err
}
