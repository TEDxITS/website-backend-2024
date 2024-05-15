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
	PreEvent3Controller interface {
		RegisterPreEvent3(ctx *gin.Context)
		GetPreEvent3Status(ctx *gin.Context)
	}

	preEvent3Controller struct {
		preEvent3Service service.PreEvent3Service
	}
)

func NewPreEvent3Controller(service service.PreEvent3Service) PreEvent3Controller {
	return &preEvent3Controller{
		preEvent3Service: service,
	}
}

func (c *preEvent3Controller) RegisterPreEvent3(ctx *gin.Context) {
	var req dto.PE3RSVPRegister
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := c.preEvent3Service.RegisterPE3(req, ctx.GetString(constants.CTX_KEY_USER_ID))
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TICKET, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_TICKET, nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *preEvent3Controller) GetPreEvent3Status(ctx *gin.Context) {
	status, err := c.preEvent3Service.GetStatus()
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_PE3_STATUS, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_PE3_STATUS, status)
	ctx.JSON(http.StatusOK, res)
}
