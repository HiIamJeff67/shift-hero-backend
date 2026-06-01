package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/shift-hero-backend/app/dtos"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
)

type CompanyControllerInterface interface {
	CreateCompany(ctx *gin.Context, reqDto *dtos.CreateCompanyReqDto)
	GetMyCompanies(ctx *gin.Context, reqDto *dtos.GetMyCompaniesReqDto)
	GetCompany(ctx *gin.Context, reqDto *dtos.GetCompanyReqDto)
	UpdateCompany(ctx *gin.Context, reqDto *dtos.UpdateCompanyReqDto)
	GetCompanyMembers(ctx *gin.Context, reqDto *dtos.GetCompanyMembersReqDto)
	AddCompanyMember(ctx *gin.Context, reqDto *dtos.AddCompanyMemberReqDto)
	UpdateCompanyMember(ctx *gin.Context, reqDto *dtos.UpdateCompanyMemberReqDto)
	DeleteCompanyMember(ctx *gin.Context, reqDto *dtos.DeleteCompanyMemberReqDto)
}

type CompanyController struct {
	companyService services.CompanyServiceInterface
}

func NewCompanyController(service services.CompanyServiceInterface) CompanyControllerInterface {
	return &CompanyController{companyService: service}
}

func (c *CompanyController) CreateCompany(ctx *gin.Context, reqDto *dtos.CreateCompanyReqDto) {
	resDto, exception := c.companyService.CreateCompany(ctx.Request.Context(), reqDto)
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

func (c *CompanyController) GetMyCompanies(ctx *gin.Context, reqDto *dtos.GetMyCompaniesReqDto) {
	resDto, exception := c.companyService.GetMyCompanies(ctx.Request.Context(), reqDto)
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

func (c *CompanyController) GetCompany(ctx *gin.Context, reqDto *dtos.GetCompanyReqDto) {
	resDto, exception := c.companyService.GetCompany(ctx.Request.Context(), reqDto)
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

func (c *CompanyController) UpdateCompany(ctx *gin.Context, reqDto *dtos.UpdateCompanyReqDto) {
	resDto, exception := c.companyService.UpdateCompany(ctx.Request.Context(), reqDto)
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

func (c *CompanyController) GetCompanyMembers(ctx *gin.Context, reqDto *dtos.GetCompanyMembersReqDto) {
	resDto, exception := c.companyService.GetCompanyMembers(ctx.Request.Context(), reqDto)
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

func (c *CompanyController) AddCompanyMember(ctx *gin.Context, reqDto *dtos.AddCompanyMemberReqDto) {
	resDto, exception := c.companyService.AddCompanyMember(ctx.Request.Context(), reqDto)
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

func (c *CompanyController) UpdateCompanyMember(ctx *gin.Context, reqDto *dtos.UpdateCompanyMemberReqDto) {
	resDto, exception := c.companyService.UpdateCompanyMember(ctx.Request.Context(), reqDto)
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

func (c *CompanyController) DeleteCompanyMember(ctx *gin.Context, reqDto *dtos.DeleteCompanyMemberReqDto) {
	resDto, exception := c.companyService.DeleteCompanyMember(ctx.Request.Context(), reqDto)
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
