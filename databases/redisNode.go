/*
 * Author: Yang Aobo
 * Telegram: @AnAsianGangster
 * Created At: March 17, 2021
 * Updated At: March 17, 2021
 * Last Modified By: Yang Aobo
 */

/**
 * This package contains functions do database I/O operations
 *
 *
 * This file contains functions that create & delete on redis HSET
 */
package databases

import (
	"fmt"

	"github.com/go-redis/redis"
)

func CreateNode(node string, redisClient *redis.Client, firstKey string, firstValue string) bool {
	redisResponse, err := redisClient.Do("HSET", node, firstKey, firstValue).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("HSET failed, creating node failed")
			return false
		}
		panic(err)
	}

	if redisResponse.(int64) == 1 {
		return true
	} else {
		return false
	}
}

func DeleteNode(node string, redisClient *redis.Client) bool {
	redisResponse, err := redisClient.Do("DEL", node).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("DEL failed, deleting node failed")
			return false
		}
		panic(err)
	}

	if redisResponse.(int64) == 1 {
		return true
	} else {
		return false
	}
}
