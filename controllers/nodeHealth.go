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

	"github.com/gin-gonic/gin"
)

// TODO change to enums & move to a private folder
var NodeStatus string = "ALIVE"

// TODO I use Content-Type: application/json. Might need to change to postform
type nodeStatusStruct struct {
	Status string `json:"status"`
}

func CheckNodeStatus(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": NodeStatus,
	})
}

func ChangeNodeStatus(context *gin.Context) {
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
