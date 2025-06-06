package utils

import (
	"context"
	"os"
	"time"

	"tic-tac-toe/backend/logger"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	if _, err := Rdb.Ping(Ctx).Result(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}
}
