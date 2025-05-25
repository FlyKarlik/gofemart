package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/FlyKarlik/gofemart/internal/model"
	"github.com/FlyKarlik/gofemart/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type cacheItem struct {
	User      *model.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

type UserCache struct {
	logger logger.Logger
	client *redis.Client
}

func NewUserCache(logger logger.Logger, client *redis.Client) *UserCache {
	return &UserCache{
		logger: logger,
		client: client,
	}
}

func (c *UserCache) Set(ctx context.Context, userID uuid.UUID, user *model.User, ttl time.Duration) error {
	item := cacheItem{
		User:      user,
		ExpiresAt: time.Now().Add(ttl),
	}

	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, userID.String(), jsonData, ttl).Err()
}

func (c *UserCache) Get(ctx context.Context, userID uuid.UUID) (*model.User, bool, error) {
	val, err := c.client.Get(ctx, userID.String()).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	var item cacheItem
	if err := json.Unmarshal([]byte(val), &item); err != nil {
		return nil, false, err
	}

	if time.Now().After(item.ExpiresAt) {
		go c.client.Del(context.Background(), userID.String())
		return nil, false, nil
	}

	return item.User, true, nil
}

func (c *UserCache) Delete(ctx context.Context, userID uuid.UUID) error {
	return c.client.Del(ctx, userID.String()).Err()
}
