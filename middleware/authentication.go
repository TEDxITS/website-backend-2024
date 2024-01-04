package middleware

import (
	"net/http"
	"strings"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			abortTokenInvalid(ctx)
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			abortTokenInvalid(ctx)
			return
		}

		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			var response utils.Response
			if err.Error() == dto.ErrTokenExpired.Error() {
				response = utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_TOKEN, dto.ErrTokenExpired.Error(), nil)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}
			abortTokenInvalid(ctx)
		}

		if !token.Valid {
			abortTokenInvalid(ctx)
			return
		}

		userId, userRole, err := jwtService.GetPayloadInsideToken(authHeader)
		if err != nil {
			abortTokenInvalid(ctx)
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
