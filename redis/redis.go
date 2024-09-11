package redis

import "github.com/go-redis/redis/v8"

var (
	RedisClient *redis.Client
)

func InitRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return redisClient
}