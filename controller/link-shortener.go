package controller

import (
	"net/http"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

type (
	LinkShortenerController interface {
		RedirectByAlias(ctx *gin.Context)
		Create(ctx *gin.Context)
		GetAllPagination(ctx *gin.Context)
	}

	linkShortenerController struct {
		linkShortenService service.LinkShortenerService
	}
)

func NewLinkShortenerController(service service.LinkShortenerService) LinkShortenerController {
	return &linkShortenerController{
		linkShortenService: service,
	}
}

func (c *linkShortenerController) Create(ctx *gin.Context) {
	var req dto.LinkShortenerRequest

	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.linkShortenService.CreateLinkShortener(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_LINK_SHORTENER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_LINK_SHORTENER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *linkShortenerController) RedirectByAlias(ctx *gin.Context) {
	alias := ctx.Param("alias")

	result, err := c.linkShortenService.RedirectByAlias(ctx.Request.Context(), alias)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LINK_SHORTENER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LINK_SHORTENER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *linkShortenerController) GetAllPagination(ctx *gin.Context) {
	var req dto.PaginationQuery

	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.linkShortenService.GetAllPagination(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LINK_SHORTENER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_LINK_SHORTENER,
		Data:    result.Data,
		Meta:    result.PaginationMetadata,
	}
	ctx.JSON(http.StatusOK, res)
}
