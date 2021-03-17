package routers

import (
	"go-http-redis/controllers"

	"github.com/gin-gonic/gin"
)

// RESTful
func MountDatabaseIORouter(router *gin.Engine) *gin.Engine {
	router.GET("/key-value-pair", controllers.FindOneKeyValuePair)
	router.POST("/key-value-pair", controllers.CreatOneKeyValuePair)
	router.PUT("/key-value-pair", controllers.UpdateOneKeyValuePair)
	router.DELETE("/key-value-pair", controllers.DeleteOneKeyValuePair)
	return router
}
