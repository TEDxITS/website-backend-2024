package routes

import (
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/gin-gonic/gin"
)

func Event(route *gin.Engine, eventController controller.EventController, jwtService service.JWTService) {
	routes := route.Group("/api/event")
	{
		routes.GET("/", eventController.FindAllUser)
		routes.GET("/:id", eventController.FindByIDUser)

		routes.GET("/admin", middleware.Authenticate(jwtService), eventController.FindAllAdmin)
		routes.GET("/:id/admin", middleware.Authenticate(jwtService), eventController.FindByIDAdmin)
	}
}
