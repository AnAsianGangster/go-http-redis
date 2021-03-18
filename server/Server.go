package server

import (
	"fmt"

	"github.com/go-redis/redis"
)

func CreateServer(serverName string, redisClient *redis.Client) bool {
	redisResponse, err := redisClient.Do("HSET", serverName, "firstKey", "firstValue").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("HSET failed, creating server failed")
			return false
		}
		panic(err)
	}

	if redisResponse.(int64) == 0 {
		fmt.Printf("%v --- has been succesfully created", serverName)
		return true
	} else {
		fmt.Printf("%v --- failed to be created", serverName)
		return false
	}
}

func AddKeyValuePair(serverName string, redisClient *redis.Client, key string, value string) bool {
	redisResponse, err := redisClient.Do("HSET", serverName, key, value).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("HSET failed, add key value pair failed")
			return false
		}
		panic(err)
	}

	if redisResponse.(int64) == 0 {
		return true
	} else {
		return false
	}
}

func GetAllKeyValuePair(serverName string, redisClient *redis.Client) []interface{} {
	redisResponse, err := redisClient.Do("HGETALL", serverName).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("HGETALL failed")
			return nil
		}
		panic(err)
	}
	return redisResponse.([]interface{})
}

func FindOneKeyValuePair(serverName string, redisClient *redis.Client, key string) string {
	redisResponse, err := redisClient.Do("HGET", serverName, key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("HGET failed, creating server failed")
			return "" // FIXME return nil
		}
		panic(err)
	}

	return redisResponse.(string)
}
