package routers

import (
	"go-http-redis/controllers"

	"github.com/gin-gonic/gin"
)

// RESTful
func ServerOperation(router *gin.Engine) *gin.Engine {
	router.POST("/node", controllers.CreateOneServer)
	router.DELETE("/node", controllers.DeleteOneServer)
	return router
}
