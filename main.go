/*
 * Author: Yang Aobo
 * Telegram: @AnAsianGangster
 * Created At: March 17, 2021
 * Updated At: April 7, 2021
 * Last Modified By: Yang Aobo
 */

/**
 * This is the main package, the entry point
 */
package main

import (
	"go-http-redis/routers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// create the router
	router := gin.Default()

	// mount router to routes
	routers.MountDatabaseIORouter(router)
	routers.ServerOperation(router)
	routers.NodeHealth(router)

	// check if running server on go-http-redis independently or under docker
	var PORT string
	if os.Getenv("SERVER_PORT") == "" {
		PORT = "5000"
	} else {
		PORT = os.Getenv("SERVER_PORT")
	}

	router.Run(":" + PORT)
}
