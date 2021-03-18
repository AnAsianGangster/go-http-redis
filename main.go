package main

import (
	"go-http-redis/config"
	"go-http-redis/routers"
	"go-http-redis/server"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func runNode(router *gin.Engine) {
	router.Run(":5000")
}

func main() {

	// connection to the redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisConfig().ADDRESS,
		Password: config.GetRedisConfig().PASSWORD,
		DB:       int(config.GetRedisConfig().DB),
	})

	// create server:1
	createServer1Result := server.CreateServer("server:1", redisClient)
	if createServer1Result == false {
		return
	}

	router := gin.Default()
	// mount router to routes
	routers.MountDatabaseIORouter(router)
	routers.ServerOperation(router)
	// handle routing
	router.POST("/get", func(c *gin.Context) {
		key := c.PostForm("key")

		// value := server.FindOneKeyValuePair("server:1", redisClient, key)

		c.JSON(200, gin.H{
			"status": "fetched",
			"key":    key,
			// "value":  value,
		})
	})

	router.POST("/set", func(c *gin.Context) {
		// key := c.PostForm("key")

		// value := server.FindOneKeyValuePair("server:1", redisClient, key)

		c.JSON(200, gin.H{
			"status": "posted",
			// "value":  value,
		})
	})

	// running on port
	go runNode(router) // run goroutine to run the server and prevent blocking

	for {
		select {
		default:
			// do something here
			// create the router

		}
	}
}
