package binders

import (
	"github.com/gin-gonic/gin"

	contexts "github.com/your-org/go-start-monolithic-kit/app/contexts"
	dtos "github.com/your-org/go-start-monolithic-kit/app/dtos"
	types "github.com/your-org/go-start-monolithic-kit/shared/types"
)

type UserSettingBinderInterface interface {
	BindGetMySetting(controllerFunc types.ControllerFunc[*dtos.GetMySettingReqDto]) gin.HandlerFunc
}

type UserSettingBinder struct{}

func NewUserSettingBinder() UserSettingBinderInterface {
	return &UserSettingBinder{}
}

func (b *UserSettingBinder) BindGetMySetting(controllerFunc types.ControllerFunc[*dtos.GetMySettingReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMySettingReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		controllerFunc(ctx, &reqDto)
	}
}
