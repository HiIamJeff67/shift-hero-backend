package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/your-org/go-start-monolithic-kit/app/dtos"
	services "github.com/your-org/go-start-monolithic-kit/app/services"
)

type UserSettingControllerInterface interface {
	GetMySetting(ctx *gin.Context, reqDto *dtos.GetMySettingReqDto)
}

type UserSettingController struct {
	userSettingService services.UserSettingServiceInterface
}

func NewUserSettingController(service services.UserSettingServiceInterface) UserSettingControllerInterface {
	return &UserSettingController{
		userSettingService: service,
	}
}

func (c *UserSettingController) GetMySetting(ctx *gin.Context, reqDto *dtos.GetMySettingReqDto) {
	resDto, exception := c.userSettingService.GetMySetting(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
