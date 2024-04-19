package routes

import (
	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/gin-gonic/gin"
)

func Storage(route *gin.Engine, storageController controller.StorageController, jwtService config.JWTService) {
	routes := route.Group("/api/storage")
	{
		routes.GET("/main-event/:id", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), storageController.GetMainEventPaymentFile)
	}
}
