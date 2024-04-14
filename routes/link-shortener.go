package routes

import (
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/controller"
	"github.com/TEDxITS/website-backend-2024/middleware"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/gin-gonic/gin"
)

func LinkShortener(route *gin.Engine, linkShortenerController controller.LinkShortenerController, jwtService service.JWTService) {
	routes := route.Group("/api/links")
	{
		routes.GET("/:alias", linkShortenerController.RedirectByAlias)
		routes.GET("", middleware.Authenticate(jwtService), middleware.OnlyAllow(constants.ENUM_ROLE_ADMIN), linkShortenerController.GetAllPagination)
		routes.POST("", linkShortenerController.Create)
	}
}
