package tools

import (
	"go-http-redis/config"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func init() {
	// connection to the redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.GetRedisConfig().ADDRESS,
		Password: config.GetRedisConfig().PASSWORD,
		DB:       int(config.GetRedisConfig().DB),
	})
}

//GetConfig - get singleton instance pre-initialized
func GetRedisClient() *redis.Client {
	return redisClient
}
