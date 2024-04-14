package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/service"
)

func PreEvent2(route *gin.Engine, ticketController controller.PreEvent2Controller, jwtService service.JWTService) {
	routes := route.Group("/api/ticket")
	{
		routes.POST("/pre-event-2", ticketController.CreatePE2RSVP)
		routes.GET("/pre-event-2", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), ticketController.GetPE2RSVPPaginated)
		routes.GET("/pre-event-2/counter", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), ticketController.GetPE2RSVPCounter)
		routes.GET("/pre-event-2/status", ticketController.GetPE2RSVPStatus)
		routes.GET("/pre-event-2/:id", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), ticketController.GetPE2RSVPDetail)

		routes.POST("/main-event/confirm-payment", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), ticketController.ConfirmPaymentME)
		routes.POST("/main-event/check-in", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), ticketController.CheckInME)
		routes.GET("/main-event/status", middleware.Authenticate(jwtService), ticketController.GetMEStatus)
	}
}