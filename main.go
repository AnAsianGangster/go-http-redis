package main

import (
	"go-http-redis/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// create the router
	router := gin.Default()

	// mount router to routes
	routers.MountDatabaseIORouter(router)
	routers.ServerOperation(router)

	router.Run(":5000")
}
