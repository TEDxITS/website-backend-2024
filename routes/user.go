package routes

import (
	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/gin-gonic/gin"
)

func User(route *gin.Engine, userController controller.UserController, jwtService config.JWTService) {
	routes := route.Group("/api/user")
	{
		routes.POST("", userController.Register)
		routes.POST("/login", userController.Login)
		routes.PATCH("", middleware.Authenticate(jwtService), userController.Update)
		routes.GET("/me", middleware.Authenticate(jwtService), userController.Me)
		routes.GET("", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), userController.GetAllPagination)
		routes.GET("/verify", userController.Verify)
		routes.POST("/verify/resend", userController.ResendVerifyEmail)
		routes.POST("/send-reset-password", userController.SendResetPasswordEmail)
		routes.POST("/reset-password", userController.ResetPassword)
	}
}
