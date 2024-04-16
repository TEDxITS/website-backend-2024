package routes

import (
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/gin-gonic/gin"
)

func MainEvent(route *gin.Engine, mainEventController controller.MainEventController, jwtService service.JWTService) {
	routes := route.Group("/api/ticket")
	{
		routes.POST("/main-event", middleware.Authenticate(jwtService))
		routes.POST("/main-event/check-in", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), mainEventController.CheckIn)
		routes.POST("/main-event/confirm-payment", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), mainEventController.ConfirmPayment)
		routes.GET("/main-event", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), mainEventController.GetMainEventPaginated)
		routes.GET("/main-event/counter", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), mainEventController.GetMainEventCounter)
		routes.GET("/main-event/status", mainEventController.GetStatus)
		routes.GET("/main-event/status/early-bird")
		routes.GET("/main-event/status/pre-sale")
		routes.GET("/main-event/status/normal")
		routes.GET("/main-event/:id", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), mainEventController.GetMainEventDetail)
	}
}
