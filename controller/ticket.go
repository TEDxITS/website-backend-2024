package controller

import (
	"net/http"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"

	"github.com/gin-gonic/gin"
)

type (
	TicketController interface {
		CreatePE2RSVP(ctx *gin.Context)
		GetPE2RSVPPaginated(ctx *gin.Context)
		GetPE2RSVPDetail(ctx *gin.Context)
	}

	ticketController struct {
		ticketService service.TicketService
	}
)

func NewTicketController(service service.TicketService) TicketController {
	return &ticketController{
		ticketService: service,
	}
}

func (c *ticketController) CreatePE2RSVP(ctx *gin.Context) {
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

func (c *ticketController) GetPE2RSVPPaginated(ctx *gin.Context) {
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

func (c *ticketController) GetPE2RSVPDetail(ctx *gin.Context) {
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