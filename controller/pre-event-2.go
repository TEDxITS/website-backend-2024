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
		GetMainEventPaginated(ctx *gin.Context)
		GetMainEventDetail(ctx *gin.Context)
		GetMainEventCounter(ctx *gin.Context)
	}

	preEvent2Controller struct {
		preevent2Service service.PreEvent2Service
	}
)

func NewPreEvent2Controller(service service.PreEvent2Service) PreEvent2Controller {
	return &preEvent2Controller{
		preevent2Service: service,
	}
}

func (c *preEvent2Controller) CreatePE2RSVP(ctx *gin.Context) {
	var req dto.PE2RSVPRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.preevent2Service.CreatePE2RSVP(ctx.Request.Context(), req)
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

	result, err := c.preevent2Service.GetPE2RSVPPaginated(ctx.Request.Context(), req)
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

	result, err := c.preevent2Service.GetPE2RSVPDetail(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *preEvent2Controller) GetPE2RSVPCounter(ctx *gin.Context) {
	result, err := c.preevent2Service.GetPE2RSVPCounter(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *preEvent2Controller) GetPE2RSVPStatus(ctx *gin.Context) {
	result, err := c.preevent2Service.GetPE2RSVPStatus(ctx.Request.Context())
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

func (c *ticketController) GetMainEventPaginated(ctx *gin.Context) {
	var req dto.PaginationQuery
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.ticketService.GetMainEventPaginated(ctx.Request.Context(), req)
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

func (c *ticketController) GetMainEventDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.ticketService.GetMainEventDetail(ctx.Request.Context(), id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *ticketController) GetMainEventCounter(ctx *gin.Context) {
	result, err := c.ticketService.GetMainEventCounter(ctx.Request.Context())
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_TICKET, result)
	ctx.JSON(http.StatusOK, res)
}
