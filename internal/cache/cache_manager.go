package cache_manager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/razvanmarinn/chatroom/internal/cfg"
	"github.com/redis/go-redis/v9"
)

type CacheManager interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	AddToList(ctx context.Context, key string, value string) error
	GetList(ctx context.Context, key string) ([]string, error)
	DeleteFromList(ctx context.Context, key string) error
	GetLengthForList(ctx context.Context, key string) (int64, error)
}

type RedisCacheManager struct {
	DB *redis.Client
}

func (r *RedisCacheManager) Get(ctx context.Context, key string) (string, error) {
	value, err := r.DB.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New("key not found")
	} else if err != nil {
		return "", err
	}
	return value, nil
}

func (r *RedisCacheManager) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.DB.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCacheManager) AddToList(ctx context.Context, key string, value string) error {
	return r.DB.RPush(ctx, key, value).Err()
}

func (r *RedisCacheManager) DeleteFromList(ctx context.Context, key string) error {
	return r.DB.RPop(ctx, key).Err()
}

func (r *RedisCacheManager) GetList(ctx context.Context, key string) ([]string, error) {
	return r.DB.LRange(ctx, key, 0, -1).Result()
}

func (r *RedisCacheManager) GetLengthForList(ctx context.Context, key string) (int64, error) {
	return r.DB.LLen(ctx, key).Result()
}

func (r *RedisCacheManager) Init(ctx context.Context, config cfg.Config) error {
	// redis_db, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 0)
	r.DB = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		// Password: os.Getenv("REDIS_PASSWORD"),
		// DB:       int(redis_db),
	})
	_, err := r.DB.Ping(ctx).Result()

	fmt.Println("Connected to redis ")
	return err
}

type MemcachedCacheManager struct {
}

func (m *MemcachedCacheManager) Get(ctx context.Context, key string) (string, error) {
	return "", errors.New("Memcached not implemented yet")
}
func (m *MemcachedCacheManager) GetLengthForList(ctx context.Context, key string) (int64, error) {
	return 0, errors.New("Memcached not implemented yet")
}
func (m *MemcachedCacheManager) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return errors.New("Memcached not implemented yet")
}

func (m *MemcachedCacheManager) Init(ctx context.Context, config cfg.Config) error {
	return errors.New("Memcached not implemented yet")
}

func (m *MemcachedCacheManager) AddToList(ctx context.Context, key string, value string) error {
	return errors.New("Memcached list operations not implemented")
}
func (m *MemcachedCacheManager) GetList(ctx context.Context, key string) ([]string, error) {
	return nil, errors.New("memcached list operations not implemented")
}

func (m *MemcachedCacheManager) DeleteFromList(ctx context.Context, key string) error {
	return errors.New("memcached list operations not implemented")
}

func NewCacheManager(ctx context.Context, config cfg.Config) (CacheManager, error) {
	switch config.CacheType {
	case "redis":
		manager := &RedisCacheManager{}
		if err := manager.Init(ctx, config); err != nil {
			return nil, err
		}
		return manager, nil
	case "memcached":
		manager := &MemcachedCacheManager{}
		if err := manager.Init(ctx, config); err != nil {
			return nil, err
		}
		return manager, nil
	default:
		return nil, errors.New("unsupported cache type")
	}
}
