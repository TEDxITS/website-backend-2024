package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/TEDxITS/website-backend-2024/controller"
)

func Ticket(route *gin.Engine, ticketController controller.TicketController) {
	routes := route.Group("/api/ticket")
	{
		routes.POST("", ticketController.Create)
	}
}
