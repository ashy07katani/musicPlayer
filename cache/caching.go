package cache

import (
	redis "github.com/redis/go-redis/v9"
)

func InitCache() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
