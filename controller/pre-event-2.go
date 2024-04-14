package controller

import (
	"net/http"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"

	"github.com/gin-gonic/gin"
)

type (
	PreEvent2Controller interface {
		CreatePE2RSVP(ctx *gin.Context)
		GetPE2RSVPPaginated(ctx *gin.Context)
		GetPE2RSVPDetail(ctx *gin.Context)
		GetPE2RSVPCounter(ctx *gin.Context)
		GetPE2RSVPStatus(ctx *gin.Context)

		ConfirmPaymentME(ctx *gin.Context)
		CheckInME(ctx *gin.Context)
		GetMEStatus(ctx *gin.Context)
	}

	preEvent2Controller struct {
		ticketService service.PreEvent2Service
	}
)

func NewTicketController(service service.PreEvent2Service) PreEvent2Controller {
	return &preEvent2Controller{
		ticketService: service,
	}
}

func (c *preEvent2Controller) CreatePE2RSVP(ctx *gin.Context) {
	var req dto.TicketPE2RSVPRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.ticketService.CreatePE2RSVP(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_TICKET, result)
	ctx.JSON(http.StatusCreated, res)
}

func (c *preEvent2Controller) GetPE2RSVPPaginated(ctx *gin.Context) {
	var req dto.PaginationQuery
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.ticketService.GetPE2RSVPPaginated(ctx.Request.Context(), req)
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

func (c *preEvent2Controller) GetPE2RSVPDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.ticketService.GetPE2RSVPDetail(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *preEvent2Controller) GetPE2RSVPCounter(ctx *gin.Context) {
	result, err := c.ticketService.GetPE2RSVPCounter(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *preEvent2Controller) GetPE2RSVPStatus(ctx *gin.Context) {
	result, err := c.ticketService.GetPE2RSVPStatus(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	status := struct {
		Status *bool `json:"status"`
	}{
		Status: &result,
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, status)
	ctx.JSON(http.StatusOK, res)
}

func (c *ticketController) ConfirmPaymentME(ctx *gin.Context) {
	var req dto.TicketMEConfirmPaymentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.ticketService.ConfirmPaymentME(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CONFIRM_PAYMENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CONFIRM_PAYMENT, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *ticketController) CheckInME(ctx *gin.Context) {
	var req dto.TicketMECheckInRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.ticketService.CheckInME(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CHECK_IN, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CHECK_IN, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *ticketController) GetMEStatus(ctx *gin.Context) {
	result, err := c.ticketService.GetMEStatus(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_EVENT, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_EVENT, result)
	ctx.JSON(http.StatusOK, res)
}
