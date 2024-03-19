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
		Create(ctx *gin.Context)
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

func (t *ticketController) Create(ctx *gin.Context) {
	var req dto.TicketRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := t.ticketService.CreateTicket(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_TICKET, result)
	ctx.JSON(http.StatusCreated, res)
}
