/*
 * Author: Yang Aobo
 * Telegram: @AnAsianGangster
 * Created At: March 17, 2021
 * Updated At: March 17, 2021
 * Last Modified By: Yang Aobo
 */

/**
 * This package contains tooling functions OR functions return structs that
 * other packages are dependent on
 *
 *
 * This file takes in redis configuration from /config folder and expose a redis
 * database connection client/driver
 * Expose via function **GetRedisClient**
 */
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
