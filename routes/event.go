package routes

import (
	"github.com/gin-gonic/gin"
)

func Event(route *gin.Engine) {
	routes := route.Group("/api/event")
	{
		routes.GET("/")
		routes.GET("/:id")

		routes.GET("/admin")
	}
}
