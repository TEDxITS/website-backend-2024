package middleware

import (
	"fmt"
	"net/http"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

func OnlyAllow(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.GetString(constants.CTX_KEY_ROLE_NAME)

		for _, role := range roles {
			if userRole == role {
				ctx.Next()
				return
			}
		}

		err := fmt.Sprintf(dto.ErrRoleNotAllowed.Error(), userRole)
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_VERIFY_TOKEN, err, nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
}
