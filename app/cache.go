// filepath: /C:/Users/catal/Desktop/hetic/CDN-project/app/cache.go
package app

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Cache duration
const cacheDuration = 10 * time.Minute

// Function to cache data in Redis
func CacheData(key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        log.Printf("Error marshalling data for key %s: %v", key, err)
        return err
    }

    err = RedisClient.Set(context.Background(), key, data, cacheDuration).Err()
    if err != nil {
        log.Printf("Error setting cache for key %s: %v", key, err)
        return err
    }

    log.Printf("Data cached successfully for key %s", key)
    return nil
}

// Function to get cached data from Redis
func GetCachedData(key string, dest interface{}) error {
    log.Printf("Attempting to get cached data for key %s", key)
    data, err := RedisClient.Get(context.Background(), key).Result()
    if err == redis.Nil {
        log.Printf("Cache miss for key %s", key)
        return nil // Cache miss
    } else if err != nil {
        log.Printf("Error getting cached data for key %s: %v", key, err)
        return err
    }

    err = json.Unmarshal([]byte(data), dest)
    if err != nil {
        log.Printf("Error unmarshalling cached data for key %s: %v", key, err)
        return err
    }

    log.Printf("Cached data retrieved successfully for key %s", key)
    return nil
}

// Function to delete cached data from Redis
func DeleteCachedData(key string) error {
    err := RedisClient.Del(context.Background(), key).Err()
    if err != nil {
        return err
    }

    return nil
}

// Function to check if data is cached in Redis
func IsDataCached(key string) (bool, error) {
    _, err := RedisClient.Get(context.Background(), key).Result()
    if err == redis.Nil {
        return false, nil
    } else if err != nil {
        return false, err
    }

    return true, nil
}
