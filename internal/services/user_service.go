package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/razvanmarinn/chatroom/internal/logger"

	cache "github.com/razvanmarinn/chatroom/internal/cache"
	"github.com/razvanmarinn/chatroom/internal/db"
	r_fact "github.com/razvanmarinn/chatroom/internal/db/repository_factory"
)

type UserService struct {
	UserRepo     db.UserRepository
	CacheManager cache.CacheManager
	Logger       logger.Logger
}

func NewUserService(cm cache.CacheManager, rf r_fact.RepositoryFactory, logger logger.Logger) (*UserService, error) {
	userRepo, err := rf.CreateUserRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	return &UserService{
		UserRepo:     userRepo,
		CacheManager: cm,
		Logger:       logger,
	}, nil
}

func (us *UserService) GetUserByUsername(username string) (db.User, error) {
	// cachedUser, err := us.CacheManager.Get(username)
	// if err == nil && cachedUser != nil {
	// 	return *cachedUser, nil
	// }

	user, err := us.UserRepo.GetUserByUsername(username)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to retrieve user by username: %w", err)
	}

	// err = us.CacheManager.Set(username, &user)
	// if err != nil {
	// 	log.Printf("Warning: failed to cache user %s: %v", username, err)
	// }

	return user, nil
}

func (us *UserService) GetUserByID(uuid uuid.UUID) (db.User, error) {
	// cachedUser, err := us.CacheManager.Get(username)
	// if err == nil && cachedUser != nil {
	// 	return *cachedUser, nil
	// }

	user, err := us.UserRepo.GetUserByID(uuid)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to retrieve user by username: %w", err)
	}

	// err = us.CacheManager.Set(username, &user)
	// if err != nil {
	// 	log.Printf("Warning: failed to cache user %s: %v", username, err)
	// }

	return user, nil
}
