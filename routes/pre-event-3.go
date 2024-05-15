package routes

import (
	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/gin-gonic/gin"
)

func PreEvent3(r *gin.Engine, c controller.PreEvent3Controller, jwt config.JWTService) {
	preEvent3 := r.Group("/api/ticket/pre-event-3")
	{
		preEvent3.POST("", middleware.Authenticate(jwt), c.RegisterPreEvent3)
		preEvent3.GET("/status", c.GetPreEvent3Status)
	}
}
