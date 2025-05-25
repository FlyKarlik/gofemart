package database

import (
	"time"

	"github.com/FlyKarlik/gofemart/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(config *config.Redis) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Host + ":" + config.Port,
		MinIdleConns: config.MinIdleConns,
		PoolSize:     config.PoolSize,
		PoolTimeout:  time.Duration(config.PoolTimeout) * time.Second,
		Password:     "",
		DB:           0,
	})

	return client
}
