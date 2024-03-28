package controller

import (
	"net/http"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

type (
	EventController interface {
		FindAll(ctx *gin.Context)
		FindByID(ctx *gin.Context)
	}

	eventController struct {
		eventService service.EventService
	}
)

func NewEventController(es service.EventService) EventController {
	return &eventController{
		eventService: es,
	}
}

func (c *eventController) FindAll(ctx *gin.Context) {
	userRole := ctx.GetString(constants.CTX_KEY_ROLE_NAME)

	result, err := c.eventService.FindAll(ctx, userRole)
	if err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_EVENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_EVENT, result)
	ctx.JSON(http.StatusOK, response)
}

func (c *eventController) FindByID(ctx *gin.Context) {
	userRole := ctx.GetString(constants.CTX_KEY_ROLE_NAME)
	id := ctx.Param("id")

	result, err := c.eventService.FindByID(ctx, id, userRole)
	if err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_EVENT_NOT_FOUND, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_EVENT, result)
	ctx.JSON(http.StatusOK, response)
}
