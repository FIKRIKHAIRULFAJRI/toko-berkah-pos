package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return nil
	}

	log.Println("Redis connected successfully")
	return client
}
