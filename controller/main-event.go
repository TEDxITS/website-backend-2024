package controller

import (
	"net/http"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

type (
	MainEventController interface {
		ConfirmPayment(ctx *gin.Context)
		CheckIn(ctx *gin.Context)
		GetStatus(ctx *gin.Context)
		GetMainEventPaginated(ctx *gin.Context)
		GetMainEventDetail(ctx *gin.Context)
		GetMainEventCounter(ctx *gin.Context)
	}

	mainEventController struct {
		mainEventService service.MainEventService
	}
)

func NewMainEventController(service service.MainEventService) MainEventController {
	return &mainEventController{
		mainEventService: service,
	}
}

func (c *mainEventController) ConfirmPayment(ctx *gin.Context) {
	var req dto.MainEventConfirmPaymentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.mainEventService.ConfirmPayment(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CONFIRM_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CONFIRM_PAYMENT, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *mainEventController) CheckIn(ctx *gin.Context) {
	var req dto.MainEventCheckInRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.mainEventService.CheckIn(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CHECK_IN, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CHECK_IN, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *mainEventController) GetStatus(ctx *gin.Context) {
	result, err := c.mainEventService.GetStatus(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_EVENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_EVENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *mainEventController) GetMainEventPaginated(ctx *gin.Context) {
	var req dto.PaginationQuery
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.mainEventService.GetMainEventPaginated(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_TICKET,
		Data:    result.Data,
		Meta:    result.PaginationMetadata,
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *mainEventController) GetMainEventDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.mainEventService.GetMainEventDetail(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *mainEventController) GetMainEventCounter(ctx *gin.Context) {
	result, err := c.mainEventService.GetMainEventCounter(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}
