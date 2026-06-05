package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/shift-hero-backend/app/contexts"
	dtos "github.com/HiIamJeff67/shift-hero-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

type SchedulingBinderInterface interface {
	BindCreateShiftRequirement(controllerFunc types.ControllerFunc[*dtos.CreateShiftRequirementReqDto]) gin.HandlerFunc
	BindGetShiftRequirements(controllerFunc types.ControllerFunc[*dtos.GetShiftRequirementsReqDto]) gin.HandlerFunc
	BindUpdateShiftRequirement(controllerFunc types.ControllerFunc[*dtos.UpdateShiftRequirementReqDto]) gin.HandlerFunc
	BindDeleteShiftRequirement(controllerFunc types.ControllerFunc[*dtos.DeleteShiftRequirementReqDto]) gin.HandlerFunc
	BindUpsertAvailabilitySlots(controllerFunc types.ControllerFunc[*dtos.UpsertAvailabilitySlotsReqDto]) gin.HandlerFunc
	BindGetAvailabilitySlots(controllerFunc types.ControllerFunc[*dtos.GetAvailabilitySlotsReqDto]) gin.HandlerFunc
	BindDeleteAvailabilitySlot(controllerFunc types.ControllerFunc[*dtos.DeleteAvailabilitySlotReqDto]) gin.HandlerFunc
	BindGenerateAssignments(controllerFunc types.ControllerFunc[*dtos.GenerateAssignmentsReqDto]) gin.HandlerFunc
	BindReplaceAssignments(controllerFunc types.ControllerFunc[*dtos.ReplaceAssignmentsReqDto]) gin.HandlerFunc
	BindClaimAssignment(controllerFunc types.ControllerFunc[*dtos.ClaimAssignmentReqDto]) gin.HandlerFunc
	BindGetAssignments(controllerFunc types.ControllerFunc[*dtos.GetAssignmentsReqDto]) gin.HandlerFunc
	BindCreateSwapRequest(controllerFunc types.ControllerFunc[*dtos.CreateSwapRequestReqDto]) gin.HandlerFunc
	BindGetSwapRequests(controllerFunc types.ControllerFunc[*dtos.GetSwapRequestsReqDto]) gin.HandlerFunc
	BindClaimSwapRequest(controllerFunc types.ControllerFunc[*dtos.ClaimSwapRequestReqDto]) gin.HandlerFunc
	BindApproveSwapRequest(controllerFunc types.ControllerFunc[*dtos.ApproveSwapRequestReqDto]) gin.HandlerFunc
	BindCancelSwapRequest(controllerFunc types.ControllerFunc[*dtos.CancelSwapRequestReqDto]) gin.HandlerFunc
	BindGetSchedulePublication(controllerFunc types.ControllerFunc[*dtos.GetSchedulePublicationReqDto]) gin.HandlerFunc
	BindUpsertSchedulePublication(controllerFunc types.ControllerFunc[*dtos.UpsertSchedulePublicationReqDto]) gin.HandlerFunc
	BindGetCompanySettings(controllerFunc types.ControllerFunc[*dtos.GetCompanySettingsReqDto]) gin.HandlerFunc
	BindUpdateCompanySettings(controllerFunc types.ControllerFunc[*dtos.UpdateCompanySettingsReqDto]) gin.HandlerFunc
}

type SchedulingBinder struct{}

func NewSchedulingBinder() SchedulingBinderInterface {
	return &SchedulingBinder{}
}

func extractUserId(ctx *gin.Context) (uuid.UUID, *exceptions.Exception) {
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
	if exception != nil {
		return uuid.Nil, exception
	}
	return *userId, nil
}

func parseCompanyIdFromPathForSchedulingBinder(ctx *gin.Context) (uuid.UUID, *exceptions.Exception) {
	companyIdString := ctx.Param("companyId")
	companyId, err := uuid.Parse(companyIdString)
	if err != nil {
		return uuid.Nil, exceptions.Scheduling.BadRequest("Invalid companyId in path").WithOrigin(err)
	}
	return companyId, nil
}

func (b *SchedulingBinder) BindCreateShiftRequirement(controllerFunc types.ControllerFunc[*dtos.CreateShiftRequirementReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateShiftRequirementReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindGetShiftRequirements(controllerFunc types.ControllerFunc[*dtos.GetShiftRequirementsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetShiftRequirementsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		companyId, exception := parseCompanyIdFromPathForSchedulingBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		if err := ctx.ShouldBindQuery(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindUpdateShiftRequirement(controllerFunc types.ControllerFunc[*dtos.UpdateShiftRequirementReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateShiftRequirementReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindDeleteShiftRequirement(controllerFunc types.ControllerFunc[*dtos.DeleteShiftRequirementReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteShiftRequirementReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindUpsertAvailabilitySlots(controllerFunc types.ControllerFunc[*dtos.UpsertAvailabilitySlotsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpsertAvailabilitySlotsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindGetAvailabilitySlots(controllerFunc types.ControllerFunc[*dtos.GetAvailabilitySlotsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAvailabilitySlotsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		companyId, exception := parseCompanyIdFromPathForSchedulingBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		if err := ctx.ShouldBindQuery(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindDeleteAvailabilitySlot(controllerFunc types.ControllerFunc[*dtos.DeleteAvailabilitySlotReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteAvailabilitySlotReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindGenerateAssignments(controllerFunc types.ControllerFunc[*dtos.GenerateAssignmentsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GenerateAssignmentsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindReplaceAssignments(controllerFunc types.ControllerFunc[*dtos.ReplaceAssignmentsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ReplaceAssignmentsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindClaimAssignment(controllerFunc types.ControllerFunc[*dtos.ClaimAssignmentReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ClaimAssignmentReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindGetAssignments(controllerFunc types.ControllerFunc[*dtos.GetAssignmentsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAssignmentsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		companyId, exception := parseCompanyIdFromPathForSchedulingBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		if err := ctx.ShouldBindQuery(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindCreateSwapRequest(controllerFunc types.ControllerFunc[*dtos.CreateSwapRequestReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateSwapRequestReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindGetSwapRequests(controllerFunc types.ControllerFunc[*dtos.GetSwapRequestsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetSwapRequestsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		companyId, exception := parseCompanyIdFromPathForSchedulingBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		if err := ctx.ShouldBindQuery(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindClaimSwapRequest(controllerFunc types.ControllerFunc[*dtos.ClaimSwapRequestReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ClaimSwapRequestReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindApproveSwapRequest(controllerFunc types.ControllerFunc[*dtos.ApproveSwapRequestReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.ApproveSwapRequestReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindCancelSwapRequest(controllerFunc types.ControllerFunc[*dtos.CancelSwapRequestReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CancelSwapRequestReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindGetSchedulePublication(controllerFunc types.ControllerFunc[*dtos.GetSchedulePublicationReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetSchedulePublicationReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		companyId, exception := parseCompanyIdFromPathForSchedulingBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		if err := ctx.ShouldBindQuery(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindUpsertSchedulePublication(controllerFunc types.ControllerFunc[*dtos.UpsertSchedulePublicationReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpsertSchedulePublicationReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindGetCompanySettings(controllerFunc types.ControllerFunc[*dtos.GetCompanySettingsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetCompanySettingsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		companyId, exception := parseCompanyIdFromPathForSchedulingBinder(ctx)
		if exception != nil {
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.CompanyId = companyId
		controllerFunc(ctx, &reqDto)
	}
}

func (b *SchedulingBinder) BindUpdateCompanySettings(controllerFunc types.ControllerFunc[*dtos.UpdateCompanySettingsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateCompanySettingsReqDto
		uid, exception := extractUserId(ctx)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = uid
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Scheduling.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}
