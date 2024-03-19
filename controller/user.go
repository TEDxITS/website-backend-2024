package controller

import (
	"net/http"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/TEDxITS/website-backend-2024/service"
	"github.com/TEDxITS/website-backend-2024/utils"
	"github.com/gin-gonic/gin"
)

type (
	UserController interface {
		Register(ctx *gin.Context)
		Login(ctx *gin.Context)
		Update(ctx *gin.Context)
		Me(ctx *gin.Context)
		GetAllPagination(ctx *gin.Context)
	}

	userController struct {
		jwtService  service.JWTService
		userService service.UserService
	}
)

func NewUserController(us service.UserService, jwt service.JWTService) UserController {
	return &userController{
		jwtService:  jwt,
		userService: us,
	}
}

func (c *userController) Register(ctx *gin.Context) {
	var user dto.UserRequest
	if err := ctx.ShouldBind(&user); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = c.userService.SendVerificationEmail(result.Email)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// email, err := c.userService.GenerateVerificationEmail(result.Email)
	// if err != nil {
	// 	res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
	// 	ctx.JSON(http.StatusInternalServerError, res)
	// 	return
	// }

	// err = utils.SendMail(email)
	// if err != nil {
	// 	res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
	// 	ctx.JSON(http.StatusInternalServerError, res)
	// 	fmt.Println("OR HERE !?")
	// 	return
	// }

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REGISTER_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := c.userService.VerifyLogin(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	token := c.jwtService.GenerateToken(user.ID.String(), user.Role.Name)
	userResponse := entity.Authorization{
		Token: token,
		Role:  user.Role.Name,
	}

	response := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGIN, userResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) Update(ctx *gin.Context) {
	var req dto.UserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	userId := ctx.MustGet("user_id").(string)
	result, err := c.userService.UpdateUser(ctx.Request.Context(), req, userId)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Me(ctx *gin.Context) {
	userId := ctx.MustGet(constants.CTX_KEY_USER_ID).(string)
	userRole := ctx.MustGet(constants.CTX_KEY_ROLE_NAME).(string)

	result, err := c.userService.Me(ctx.Request.Context(), userId, userRole)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) GetAllPagination(ctx *gin.Context) {
	var req dto.PaginationQuery
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.userService.GetAllPagination(ctx.Request.Context(), req)
	if err != nil {
		res := utils.BuildResponseFailed(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.Response{
		Status:  true,
		Message: dto.MESSAGE_SUCCESS_GET_USER,
		Data:    result.Data,
		Meta:    result.PaginationMetadata,
	}
	ctx.JSON(http.StatusOK, res)
}
