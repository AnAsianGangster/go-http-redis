package main

import (
	"fmt"
	"go-http-redis/config"
	"go-http-redis/server"
	"go-http-redis/util"

	"github.com/go-redis/redis"
)

func main() {
	// fmt.Println("Hello World!")
	// connection to the redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisConfig().ADDRESS,
		Password: config.GetRedisConfig().PASSWORD,
		DB:       int(config.GetRedisConfig().DB),
	})

	// create server:1
	createServer1Result := server.CreateServer("server:1", redisClient)
	fmt.Println(createServer1Result)
	// create server:2
	createServer2Result := server.CreateServer("server:2", redisClient)
	fmt.Println(createServer2Result)
	// add key value pair to server:1
	addKeyValuePairResult := server.AddKeyValuePair("server:1", redisClient, "secondKey", "secondValue")
	fmt.Println(addKeyValuePairResult)
	// add key value pair to server:1

	// check server:1
	server1AllVal := server.GetAllKeyValuePair("server:1", redisClient)
	util.PrintArrayInterface(server1AllVal)

	// check server:2
	server2AllVal := server.GetAllKeyValuePair("server:2", redisClient)
	util.PrintArrayInterface(server2AllVal)

	/* 
	// flush redis
	val, err := redisClient.Do("FLUSHALL").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("FLUSHALL failed")
			return
		}
		panic(err)
	}
	fmt.Println(val.(string))
	 */
}