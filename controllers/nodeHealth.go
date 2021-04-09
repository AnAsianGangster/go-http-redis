/*
 * Author: Yang Aobo
 * Telegram: @AnAsianGangster
 * Created At: April 7, 2021
 * Updated At: April 7, 2021
 * Last Modified By: Yang Aobo
 */

/**
 * This package contains HTTP handler functions
 *
 *
 * This file contains handler functions that handle node health related operations
 *
 * All functions destructure HTTP requests, call database operations, build response
 * and reply with response
 */
package controllers

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// TODO change to enums & move to a private folder
var NodeStatus string = "ALIVE"

// TODO I use Content-Type: application/json. Might need to change to postform
type nodeStatusStruct struct {
	Status string `json:"status"`
}

func CheckNodeStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, gin.H{
			"name":     os.Getenv("SERVER_NAME"),
			"nodeName": os.Getenv("DOCKER_SERVER_NAME"),
			"port":     os.Getenv("SERVER_PORT"),
			"status":   NodeStatus,
		})
	}
}

func ChangeNodeStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		var nodeStatusStruct nodeStatusStruct
		err := context.BindJSON(&nodeStatusStruct)
		if err != nil {
			log.Fatal(err)
		}

		NodeStatus = nodeStatusStruct.Status

		context.JSON(200, gin.H{
			"success": true,
		})
	}
}
