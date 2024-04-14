package routes

import (
	"github.com/TEDxITS/website-backend-2024/websocket"
	"github.com/gin-gonic/gin"
)

func TicketQueue(route *gin.Engine, earlyBirdHandler, preSaleHandler, normalHandler websocket.TicketQueue) {
	routes := route.Group("/ws")
	{
		routes.GET("/early-bird", earlyBirdHandler.Serve)
		routes.GET("/pre-sale", preSaleHandler.Serve)
		routes.GET("/normal", normalHandler.Serve)
	}
}
