package routers

import (
	"go-http-redis/controllers"

	"github.com/gin-gonic/gin"
)

// RPC
func ServerOperation(router *gin.Engine) *gin.Engine {
	router.GET("/create-server", controllers.CreateOneServer)
	router.DELETE("/delete-server", controllers.DeleteOneServer)
	return router
}
