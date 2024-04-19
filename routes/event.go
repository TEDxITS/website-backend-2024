package routes

import (
	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/gin-gonic/gin"
)

func Event(route *gin.Engine, eventController controller.EventController, jwtService config.JWTService) {
	routes := route.Group("/api/events")
	{
		routes.GET("/", middleware.Authenticate(jwtService), eventController.FindAll)
		routes.GET("/pre-event-2", middleware.Authenticate(jwtService), eventController.GetPE2Detail)
		routes.GET("/:id", middleware.Authenticate(jwtService), eventController.FindByID)
	}
}
