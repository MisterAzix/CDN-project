package app

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)


func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestCacheData(t *testing.T) {
	key := "testKey"
	value := "testValue"

	err := CacheData(key, value)
	assert.NoError(t, err)

	// Clean up
	defer RedisClient.Del(context.Background(), key)
}

func TestGetCachedData(t *testing.T) {
	key := "testKey"
	value := "testValue"

	err := CacheData(key, value)
	assert.NoError(t, err)

	var result string
	err = GetCachedData(key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Clean up
	defer RedisClient.Del(context.Background(), key)
}

func TestDeleteCachedData(t *testing.T) {
	key := "testKey"
	value := "testValue"

	err := CacheData(key, value)
	assert.NoError(t, err)

	err = DeleteCachedData(key)
	assert.NoError(t, err)

	var result string
	err = GetCachedData(key, &result)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestIsDataCached(t *testing.T) {
	key := "testKey"
	value := "testValue"

	err := CacheData(key, value)
	assert.NoError(t, err)

	isCached, err := IsDataCached(key)
	assert.NoError(t, err)
	assert.True(t, isCached)

	// Clean up
	defer RedisClient.Del(context.Background(), key)
}

func TestCacheMiss(t *testing.T) {
	key := "nonExistentKey"

	var result string
	err := GetCachedData(key, &result)
	assert.NoError(t, err)
	assert.Empty(t, result)

	isCached, err := IsDataCached(key)
	assert.NoError(t, err)
	assert.False(t, isCached)
}