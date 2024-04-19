package middleware

import (
	"net/http"
	"strings"

	"github.com/TEDxITS/website-backend-2024/config"
	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService config.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_TOKEN, dto.ErrTokenNotFound.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			abortTokenInvalid(ctx)
			return
		}

		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		userId, userRole, err := jwtService.GetPayloadInsideToken(authHeader)
		if err != nil {
			if err.Error() == dto.ErrTokenExpired.Error() {
				response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_TOKEN, dto.ErrTokenExpired.Error(), nil)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_TOKEN, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Set(constants.CTX_KEY_TOKEN, authHeader)
		ctx.Set(constants.CTX_KEY_USER_ID, userId)
		ctx.Set(constants.CTX_KEY_ROLE_NAME, userRole)
		ctx.Next()
	}
}

func abortTokenInvalid(ctx *gin.Context) {
	response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_TOKEN, dto.ErrTokenInvalid.Error(), nil)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
