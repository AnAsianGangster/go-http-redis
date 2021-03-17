package main

import (
	"fmt"
	"go-http-redis/config"
	"go-http-redis/routers"
	"go-http-redis/server"
	"go-http-redis/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func runNode(router *gin.Engine) {

	router.Run(":5000")

}

func main() {
	// create the router
	router := gin.Default()
	// mount router to routes
	routers.MountDatabaseIORouter(router)
	routers.ServerOperation(router)
	// handle routing
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "from 5000")
	})

	// running on port
	go runNode(router) // run goroutine to run the server and prevent blocking

	// create the router
	// router1 := gin.Default()

	// // mount router to routes
	// routers.MountDatabaseIORouter(router1)
	// routers.ServerOperation(router1)

	// router1.GET("/", func(c *gin.Context){
	// 	c.String(http.StatusOK, "from 5001")
	// })
	// // running on port
	// router1.Run(":5001")
	// // create the router
	// router2 := gin.Default()

	// // mount router to routes
	// routers.MountDatabaseIORouter(router2)
	// routers.ServerOperation(router2)

	// // running on port
	// router2.Run(":5002")
	// // create the router
	// router3 := gin.Default()

	// // mount router to routes
	// routers.MountDatabaseIORouter(router3)
	// routers.ServerOperation(router3)

	// // running on port
	// router3.Run(":5003")

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

	for {
		select {
		default:
			// do something here
		}
	}
}
