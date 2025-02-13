// filepath: /C:/Users/catal/Desktop/hetic/CDN-project/app/redis.go
package app

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedisClient() {
    redisAddr := os.Getenv("REDIS_ADDR")
    redisPassword := os.Getenv("REDIS_PASSWORD")

    RedisClient = redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword,
        DB:       0,
    })

    _, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }

    log.Println("Connected to Redis!")
}