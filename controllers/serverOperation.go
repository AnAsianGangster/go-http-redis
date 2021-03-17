package controllers

import (
	"go-http-redis/databases"
	"go-http-redis/tools"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Node       string `json:"node"`
	FirstKey   string `json:"firstKey"`
	FirstValue string `json:"firstValue"`
}

func CreateOneServer(context *gin.Context) {
	redisClient := tools.GetRedisClient()
	var server Server
	err := context.BindJSON(&server)
	if err != nil {
		log.Fatal(err)
	}

	value := databases.CreateNode(server.Node, redisClient, server.FirstKey, server.FirstValue)
	if value == true {
		context.JSON(201, gin.H{
			"node":       server.Node,
			"firstKey":   server.FirstKey,
			"firstValue": server.FirstValue,
		})
	} else {
		context.JSON(400, gin.H{
			"success": false,
		})
	}
}

func DeleteOneServer(context *gin.Context) {
	redisClient := tools.GetRedisClient()
	node := context.Query("node")
	value := databases.DeleteNode(node, redisClient)

	if value == true {
		context.JSON(200, gin.H{
			"message": "delete successfully",
			"node":    node,
		})
	} else {
		context.JSON(400, gin.H{
			"success": false,
		})
	}
}
