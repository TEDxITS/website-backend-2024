package controller

import (
	"net/http"

	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

type (
	StorageController interface {
		GetMainEventPaymentFile(c *gin.Context)
	}

	storageController struct {
		storageService service.StorageService
	}
)

func NewStorageController(Service service.StorageService) StorageController {
	return &storageController{
		storageService: Service,
	}
}

func (c *storageController) GetMainEventPaymentFile(ctx *gin.Context) {
	id := ctx.Param("id")

	file, err := c.storageService.GetMainEventPaymentFile(id)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_FILE, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	ctx.Data(http.StatusOK, "application/octet-stream", file)
}
