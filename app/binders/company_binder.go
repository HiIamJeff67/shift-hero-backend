package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/shift-hero-backend/app/contexts"
	dtos "github.com/HiIamJeff67/shift-hero-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	"github.com/HiIamJeff67/shift-hero-backend/app/monitor/logs"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

type CompanyBinderInterface interface {
	BindCreateCompany(controllerFunc types.ControllerFunc[*dtos.CreateCompanyReqDto]) gin.HandlerFunc
	BindGetMyCompanies(controllerFunc types.ControllerFunc[*dtos.GetMyCompaniesReqDto]) gin.HandlerFunc
	BindGetCompany(controllerFunc types.ControllerFunc[*dtos.GetCompanyReqDto]) gin.HandlerFunc
	BindUpdateCompany(controllerFunc types.ControllerFunc[*dtos.UpdateCompanyReqDto]) gin.HandlerFunc
	BindGetCompanyMembers(controllerFunc types.ControllerFunc[*dtos.GetCompanyMembersReqDto]) gin.HandlerFunc
	BindAddCompanyMember(controllerFunc types.ControllerFunc[*dtos.AddCompanyMemberReqDto]) gin.HandlerFunc
	BindUpdateCompanyMember(controllerFunc types.ControllerFunc[*dtos.UpdateCompanyMemberReqDto]) gin.HandlerFunc
	BindDeleteCompanyMember(controllerFunc types.ControllerFunc[*dtos.DeleteCompanyMemberReqDto]) gin.HandlerFunc
	BindCreateCompanyJoinRequest(controllerFunc types.ControllerFunc[*dtos.CreateCompanyJoinRequestReqDto]) gin.HandlerFunc
	BindGetCompanyJoinRequests(controllerFunc types.ControllerFunc[*dtos.GetCompanyJoinRequestsReqDto]) gin.HandlerFunc
	BindApproveCompanyJoinRequest(controllerFunc types.ControllerFunc[*dtos.ReviewCompanyJoinRequestReqDto]) gin.HandlerFunc
	BindRejectCompanyJoinRequest(controllerFunc types.ControllerFunc[*dtos.ReviewCompanyJoinRequestReqDto]) gin.HandlerFunc
	BindGetMyCompanyJoinRequests(controllerFunc types.ControllerFunc[*dtos.GetMyCompanyJoinRequestsReqDto]) gin.HandlerFunc
}

type CompanyBinder struct{}

func NewCompanyBinder() CompanyBinderInterface {
	return &CompanyBinder{}
}

func parseCompanyIdFromPathForCompanyBinder(ctx *gin.Context) (uuid.UUID, *exceptions.Exception) {
	companyIdString := ctx.Param("companyId")
	companyId, err := uuid.Parse(companyIdString)
	if err != nil {
		return uuid.Nil, exceptions.Company.BadRequest("Invalid companyId in path").WithOrigin(err)
	}
	return companyId, nil
}

func (b *CompanyBinder) BindCreateCompany(controllerFunc types.ControllerFunc[*dtos.CreateCompanyReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateCompanyReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindGetMyCompanies(controllerFunc types.ControllerFunc[*dtos.GetMyCompaniesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyCompaniesReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindGetCompany(controllerFunc types.ControllerFunc[*dtos.GetCompanyReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logs.Info("company-binder", "BindGetCompany entered", ctx.Request.Method, ctx.Request.URL.Path)
		var reqDto dtos.GetCompanyReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		companyId, exception := parseCompanyIdFromPathForCompanyBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindUpdateCompany(controllerFunc types.ControllerFunc[*dtos.UpdateCompanyReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateCompanyReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindGetCompanyMembers(controllerFunc types.ControllerFunc[*dtos.GetCompanyMembersReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetCompanyMembersReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		companyId, exception := parseCompanyIdFromPathForCompanyBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindAddCompanyMember(controllerFunc types.ControllerFunc[*dtos.AddCompanyMemberReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.AddCompanyMemberReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindUpdateCompanyMember(controllerFunc types.ControllerFunc[*dtos.UpdateCompanyMemberReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateCompanyMemberReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindDeleteCompanyMember(controllerFunc types.ControllerFunc[*dtos.DeleteCompanyMemberReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteCompanyMemberReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindCreateCompanyJoinRequest(controllerFunc types.ControllerFunc[*dtos.CreateCompanyJoinRequestReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateCompanyJoinRequestReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindGetCompanyJoinRequests(controllerFunc types.ControllerFunc[*dtos.GetCompanyJoinRequestsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetCompanyJoinRequestsReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		companyId, exception := parseCompanyIdFromPathForCompanyBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		if err := ctx.ShouldBindQuery(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindApproveCompanyJoinRequest(controllerFunc types.ControllerFunc[*dtos.ReviewCompanyJoinRequestReqDto]) gin.HandlerFunc {
	return b.bindReviewCompanyJoinRequest(controllerFunc)
}

func (b *CompanyBinder) BindRejectCompanyJoinRequest(controllerFunc types.ControllerFunc[*dtos.ReviewCompanyJoinRequestReqDto]) gin.HandlerFunc {
	return b.bindReviewCompanyJoinRequest(controllerFunc)
}

func (b *CompanyBinder) bindReviewCompanyJoinRequest(controllerFunc types.ControllerFunc[*dtos.ReviewCompanyJoinRequestReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ReviewCompanyJoinRequestReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Company.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *CompanyBinder) BindGetMyCompanyJoinRequests(controllerFunc types.ControllerFunc[*dtos.GetMyCompanyJoinRequestsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyCompanyJoinRequestsReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		controllerFunc(ctx, &reqDto)
	}
}
