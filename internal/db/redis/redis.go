package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Init(ctx context.Context) error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})

}
