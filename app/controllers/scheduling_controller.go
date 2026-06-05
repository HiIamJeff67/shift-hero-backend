package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/shift-hero-backend/app/dtos"
	services "github.com/HiIamJeff67/shift-hero-backend/app/services"
)

type SchedulingControllerInterface interface {
	CreateShiftRequirement(ctx *gin.Context, reqDto *dtos.CreateShiftRequirementReqDto)
	GetShiftRequirements(ctx *gin.Context, reqDto *dtos.GetShiftRequirementsReqDto)
	UpdateShiftRequirement(ctx *gin.Context, reqDto *dtos.UpdateShiftRequirementReqDto)
	DeleteShiftRequirement(ctx *gin.Context, reqDto *dtos.DeleteShiftRequirementReqDto)
	UpsertAvailabilitySlots(ctx *gin.Context, reqDto *dtos.UpsertAvailabilitySlotsReqDto)
	GetAvailabilitySlots(ctx *gin.Context, reqDto *dtos.GetAvailabilitySlotsReqDto)
	DeleteAvailabilitySlot(ctx *gin.Context, reqDto *dtos.DeleteAvailabilitySlotReqDto)
	GenerateAssignments(ctx *gin.Context, reqDto *dtos.GenerateAssignmentsReqDto)
	ReplaceAssignments(ctx *gin.Context, reqDto *dtos.ReplaceAssignmentsReqDto)
	ClaimAssignment(ctx *gin.Context, reqDto *dtos.ClaimAssignmentReqDto)
	GetAssignments(ctx *gin.Context, reqDto *dtos.GetAssignmentsReqDto)
	CreateSwapRequest(ctx *gin.Context, reqDto *dtos.CreateSwapRequestReqDto)
	GetSwapRequests(ctx *gin.Context, reqDto *dtos.GetSwapRequestsReqDto)
	ClaimSwapRequest(ctx *gin.Context, reqDto *dtos.ClaimSwapRequestReqDto)
	ApproveSwapRequest(ctx *gin.Context, reqDto *dtos.ApproveSwapRequestReqDto)
	CancelSwapRequest(ctx *gin.Context, reqDto *dtos.CancelSwapRequestReqDto)
	GetSchedulePublication(ctx *gin.Context, reqDto *dtos.GetSchedulePublicationReqDto)
	UpsertSchedulePublication(ctx *gin.Context, reqDto *dtos.UpsertSchedulePublicationReqDto)
	GetCompanySettings(ctx *gin.Context, reqDto *dtos.GetCompanySettingsReqDto)
	UpdateCompanySettings(ctx *gin.Context, reqDto *dtos.UpdateCompanySettingsReqDto)
}

type SchedulingController struct {
	schedulingService services.SchedulingServiceInterface
}

func NewSchedulingController(service services.SchedulingServiceInterface) SchedulingControllerInterface {
	return &SchedulingController{schedulingService: service}
}

func (c *SchedulingController) CreateShiftRequirement(ctx *gin.Context, reqDto *dtos.CreateShiftRequirementReqDto) {
	resDto, exception := c.schedulingService.CreateShiftRequirement(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) GetShiftRequirements(ctx *gin.Context, reqDto *dtos.GetShiftRequirementsReqDto) {
	resDto, exception := c.schedulingService.GetShiftRequirements(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) UpdateShiftRequirement(ctx *gin.Context, reqDto *dtos.UpdateShiftRequirementReqDto) {
	resDto, exception := c.schedulingService.UpdateShiftRequirement(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) DeleteShiftRequirement(ctx *gin.Context, reqDto *dtos.DeleteShiftRequirementReqDto) {
	resDto, exception := c.schedulingService.DeleteShiftRequirement(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) UpsertAvailabilitySlots(ctx *gin.Context, reqDto *dtos.UpsertAvailabilitySlotsReqDto) {
	resDto, exception := c.schedulingService.UpsertAvailabilitySlots(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) GetAvailabilitySlots(ctx *gin.Context, reqDto *dtos.GetAvailabilitySlotsReqDto) {
	resDto, exception := c.schedulingService.GetAvailabilitySlots(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) DeleteAvailabilitySlot(ctx *gin.Context, reqDto *dtos.DeleteAvailabilitySlotReqDto) {
	resDto, exception := c.schedulingService.DeleteAvailabilitySlot(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) GenerateAssignments(ctx *gin.Context, reqDto *dtos.GenerateAssignmentsReqDto) {
	resDto, exception := c.schedulingService.GenerateAssignments(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) ReplaceAssignments(ctx *gin.Context, reqDto *dtos.ReplaceAssignmentsReqDto) {
	resDto, exception := c.schedulingService.ReplaceAssignments(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) ClaimAssignment(ctx *gin.Context, reqDto *dtos.ClaimAssignmentReqDto) {
	resDto, exception := c.schedulingService.ClaimAssignment(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) GetAssignments(ctx *gin.Context, reqDto *dtos.GetAssignmentsReqDto) {
	resDto, exception := c.schedulingService.GetAssignments(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) CreateSwapRequest(ctx *gin.Context, reqDto *dtos.CreateSwapRequestReqDto) {
	resDto, exception := c.schedulingService.CreateSwapRequest(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) GetSwapRequests(ctx *gin.Context, reqDto *dtos.GetSwapRequestsReqDto) {
	resDto, exception := c.schedulingService.GetSwapRequests(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) ClaimSwapRequest(ctx *gin.Context, reqDto *dtos.ClaimSwapRequestReqDto) {
	resDto, exception := c.schedulingService.ClaimSwapRequest(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) ApproveSwapRequest(ctx *gin.Context, reqDto *dtos.ApproveSwapRequestReqDto) {
	resDto, exception := c.schedulingService.ApproveSwapRequest(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) CancelSwapRequest(ctx *gin.Context, reqDto *dtos.CancelSwapRequestReqDto) {
	resDto, exception := c.schedulingService.CancelSwapRequest(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) GetSchedulePublication(ctx *gin.Context, reqDto *dtos.GetSchedulePublicationReqDto) {
	resDto, exception := c.schedulingService.GetSchedulePublication(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) UpsertSchedulePublication(ctx *gin.Context, reqDto *dtos.UpsertSchedulePublicationReqDto) {
	resDto, exception := c.schedulingService.UpsertSchedulePublication(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) GetCompanySettings(ctx *gin.Context, reqDto *dtos.GetCompanySettingsReqDto) {
	resDto, exception := c.schedulingService.GetCompanySettings(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}

func (c *SchedulingController) UpdateCompanySettings(ctx *gin.Context, reqDto *dtos.UpdateCompanySettingsReqDto) {
	resDto, exception := c.schedulingService.UpdateCompanySettings(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": resDto, "exception": nil})
}
