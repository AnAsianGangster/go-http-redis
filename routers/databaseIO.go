/*
 * Author: Yang Aobo
 * Telegram: @AnAsianGangster
 * Created At: March 17, 2021
 * Updated At: March 17, 2021
 * Last Modified By: Yang Aobo
 */

/**
 * This package contains functions mount HTTP handlers & middlewares to URL routes
 *
 *
 * This file contains a function mount handlers to "/key-value-pair" route
 */
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
