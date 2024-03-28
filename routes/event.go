package routes

import (
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/gin-gonic/gin"
)

func Event(route *gin.Engine, eventController controller.EventController, jwtService service.JWTService) {
	routes := route.Group("/api/events")
	{
		routes.GET("/", middleware.Authenticate(jwtService), eventController.FindAll)
		routes.GET("/:id", middleware.Authenticate(jwtService), eventController.FindByID)
	}
}
