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
 * This file contains a function mount handlers to "/node" route
 */
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
