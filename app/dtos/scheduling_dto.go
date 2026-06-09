package dtos

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	"github.com/google/uuid"
)

/* ============================== Scheduling Request DTO ============================== */

type CreateShiftRequirementReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId     uuid.UUID          `json:"companyId" validate:"required,uuid4"`
			EmployeeRole  enums.EmployeeRole `json:"employeeRole" validate:"required,isemployeerole"`
			StartAt       time.Time          `json:"startAt" validate:"required"`
			EndAt         time.Time          `json:"endAt" validate:"required"`
			RequiredCount int32              `json:"requiredCount" validate:"required,min=1,max=100"`
			Note          string             `json:"note" validate:"omitempty,max=1024"`
		},
		any,
	]
}

type GetShiftRequirementsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			StartAt *time.Time `form:"startAt"`
			EndAt   *time.Time `form:"endAt"`
		},
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type UpdateShiftRequirementReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId          uuid.UUID `json:"companyId" validate:"required,uuid4"`
			ShiftRequirementId uuid.UUID `json:"shiftRequirementId" validate:"required,uuid4"`
			PartialUpdateDto[struct {
				EmployeeRole  *enums.EmployeeRole `json:"employeeRole" validate:"omitnil,isemployeerole"`
				StartAt       *time.Time          `json:"startAt" validate:"omitnil"`
				EndAt         *time.Time          `json:"endAt" validate:"omitnil"`
				RequiredCount *int32              `json:"requiredCount" validate:"omitnil,min=1,max=100"`
				Note          *string             `json:"note" validate:"omitnil,max=1024"`
			}]
		},
		any,
	]
}

type DeleteShiftRequirementReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId          uuid.UUID `json:"companyId" validate:"required,uuid4"`
			ShiftRequirementId uuid.UUID `json:"shiftRequirementId" validate:"required,uuid4"`
		},
		any,
	]
}

type UpsertAvailabilitySlotsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId uuid.UUID `json:"companyId" validate:"required,uuid4"`
			Slots     []struct {
				StartAt     time.Time `json:"startAt" validate:"required"`
				EndAt       time.Time `json:"endAt" validate:"required"`
				IsAvailable bool      `json:"isAvailable"`
			} `json:"slots" validate:"required,min=1,max=500,dive"`
		},
		any,
	]
}

type GetAvailabilitySlotsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			UserId  *uuid.UUID `form:"userId" validate:"omitnil,uuid4"`
			StartAt *time.Time `form:"startAt"`
			EndAt   *time.Time `form:"endAt"`
		},
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type DeleteAvailabilitySlotReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId          uuid.UUID `json:"companyId" validate:"required,uuid4"`
			AvailabilitySlotId uuid.UUID `json:"availabilitySlotId" validate:"required,uuid4"`
		},
		any,
	]
}

type GenerateAssignmentsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId uuid.UUID  `json:"companyId" validate:"required,uuid4"`
			StartAt   *time.Time `json:"startAt"`
			EndAt     *time.Time `json:"endAt"`
		},
		any,
	]
}

type ReplaceAssignmentsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId   uuid.UUID `json:"companyId" validate:"required,uuid4"`
			Assignments []struct {
				ShiftRequirementId uuid.UUID `json:"shiftRequirementId" validate:"required,uuid4"`
				UserId             uuid.UUID `json:"userId" validate:"required,uuid4"`
				StartAt            time.Time `json:"startAt" validate:"required"`
				EndAt              time.Time `json:"endAt" validate:"required"`
			} `json:"assignments" validate:"required,max=1000,dive"`
		},
		any,
	]
}

type ClaimAssignmentReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId          uuid.UUID `json:"companyId" validate:"required,uuid4"`
			ShiftRequirementId uuid.UUID `json:"shiftRequirementId" validate:"required,uuid4"`
		},
		any,
	]
}

type GetAssignmentsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			UserId  *uuid.UUID `form:"userId" validate:"omitnil,uuid4"`
			StartAt *time.Time `form:"startAt"`
			EndAt   *time.Time `form:"endAt"`
		},
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type CreateSwapRequestReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId         uuid.UUID `json:"companyId" validate:"required,uuid4"`
			ShiftAssignmentId uuid.UUID `json:"shiftAssignmentId" validate:"required,uuid4"`
			Reason            string    `json:"reason" validate:"omitempty,max=1024"`
		},
		any,
	]
}

type GetSwapRequestsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			Status *enums.SwapRequestStatus `form:"status" validate:"omitnil,isswaprequeststatus"`
		},
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type ClaimSwapRequestReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId     uuid.UUID `json:"companyId" validate:"required,uuid4"`
			SwapRequestId uuid.UUID `json:"swapRequestId" validate:"required,uuid4"`
		},
		any,
	]
}

type ApproveSwapRequestReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId     uuid.UUID `json:"companyId" validate:"required,uuid4"`
			SwapRequestId uuid.UUID `json:"swapRequestId" validate:"required,uuid4"`
		},
		any,
	]
}

type CancelSwapRequestReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId     uuid.UUID `json:"companyId" validate:"required,uuid4"`
			SwapRequestId uuid.UUID `json:"swapRequestId" validate:"required,uuid4"`
		},
		any,
	]
}

type GetSchedulePublicationReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			WeekStart string `form:"weekStart" validate:"required"`
			Timezone  string `form:"timezone" validate:"omitempty,max=64"`
		},
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type UpsertSchedulePublicationReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId uuid.UUID                       `json:"companyId" validate:"required,uuid4"`
			WeekStart string                          `json:"weekStart" validate:"required"`
			Timezone  string                          `json:"timezone" validate:"omitempty,max=64"`
			Status    enums.SchedulePublicationStatus `json:"status" validate:"required,isschedulepublicationstatus"`
		},
		any,
	]
}

type GetCompanySettingsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		any,
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

type UpdateCompanySettingsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			CompanyId uuid.UUID `json:"companyId" validate:"required,uuid4"`
			PartialUpdateDto[struct {
				AutoApproveSwaps *bool   `json:"autoApproveSwaps" validate:"omitnil"`
				MaxWeeklyHours   *int32  `json:"maxWeeklyHours" validate:"omitnil,min=1,max=168"`
				MinRestHours     *int32  `json:"minRestHours" validate:"omitnil,min=0,max=24"`
				Timezone         *string `json:"timezone" validate:"omitnil,max=64"`
			}]
		},
		any,
	]
}

type GenerateScheduleInsightsReqDto struct {
	Request[
		any,
		struct {
			UserId uuid.UUID
		},
		struct {
			StartAt time.Time `json:"startAt" validate:"required"`
			EndAt   time.Time `json:"endAt" validate:"required"`
			Locale  string    `json:"locale" validate:"omitempty,oneof=zh-TW en"`
			Focus   string    `json:"focus" validate:"omitempty,max=500"`
		},
		struct {
			CompanyId uuid.UUID `uri:"companyId" validate:"required,uuid4"`
		},
	]
}

/* ============================== Scheduling Response DTO ============================== */

type ShiftRequirementResDto struct {
	Id            uuid.UUID          `json:"id"`
	CompanyId     uuid.UUID          `json:"companyId"`
	EmployeeRole  enums.EmployeeRole `json:"employeeRole"`
	StartAt       time.Time          `json:"startAt"`
	EndAt         time.Time          `json:"endAt"`
	RequiredCount int32              `json:"requiredCount"`
	Note          string             `json:"note"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	CreatedAt     time.Time          `json:"createdAt"`
}

type AvailabilitySlotResDto struct {
	Id          uuid.UUID `json:"id"`
	CompanyId   uuid.UUID `json:"companyId"`
	UserId      uuid.UUID `json:"userId"`
	StartAt     time.Time `json:"startAt"`
	EndAt       time.Time `json:"endAt"`
	IsAvailable bool      `json:"isAvailable"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ShiftAssignmentResDto struct {
	Id                 uuid.UUID `json:"id"`
	CompanyId          uuid.UUID `json:"companyId"`
	ShiftRequirementId uuid.UUID `json:"shiftRequirementId"`
	UserId             uuid.UUID `json:"userId"`
	StartAt            time.Time `json:"startAt"`
	EndAt              time.Time `json:"endAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	CreatedAt          time.Time `json:"createdAt"`
}

type SwapRequestResDto struct {
	Id                uuid.UUID               `json:"id"`
	CompanyId         uuid.UUID               `json:"companyId"`
	ShiftAssignmentId uuid.UUID               `json:"shiftAssignmentId"`
	RequesterUserId   uuid.UUID               `json:"requesterUserId"`
	ClaimedByUserId   *uuid.UUID              `json:"claimedByUserId"`
	Status            enums.SwapRequestStatus `json:"status"`
	Reason            string                  `json:"reason"`
	UpdatedAt         time.Time               `json:"updatedAt"`
	CreatedAt         time.Time               `json:"createdAt"`
}

type SchedulePublicationResDto struct {
	CompanyId         uuid.UUID                       `json:"companyId"`
	WeekStart         string                          `json:"weekStart"`
	Timezone          string                          `json:"timezone"`
	Status            enums.SchedulePublicationStatus `json:"status"`
	PublishedByUserId *uuid.UUID                      `json:"publishedByUserId"`
	PublishedAt       *time.Time                      `json:"publishedAt"`
	UpdatedAt         time.Time                       `json:"updatedAt"`
	CreatedAt         time.Time                       `json:"createdAt"`
}

type CompanySettingsResDto struct {
	CompanyId        uuid.UUID `json:"companyId"`
	AutoApproveSwaps bool      `json:"autoApproveSwaps"`
	MaxWeeklyHours   int32     `json:"maxWeeklyHours"`
	MinRestHours     int32     `json:"minRestHours"`
	Timezone         string    `json:"timezone"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedAt        time.Time `json:"createdAt"`
}

type ScheduleInsightEmployeeResDto struct {
	UserId                 uuid.UUID          `json:"userId"`
	DisplayName            string             `json:"displayName"`
	EmployeeRole           enums.EmployeeRole `json:"employeeRole"`
	ShiftCount             int                `json:"shiftCount"`
	TotalHours             float64            `json:"totalHours"`
	LongestShiftHours      float64            `json:"longestShiftHours"`
	NightShiftCount        int                `json:"nightShiftCount"`
	WeekendShiftCount      int                `json:"weekendShiftCount"`
	ShortRestCount         int                `json:"shortRestCount"`
	AvailabilityConflicts  int                `json:"availabilityConflicts"`
	OpenSwapRequestCount   int                `json:"openSwapRequestCount"`
	MaxConsecutiveWorkDays int                `json:"maxConsecutiveWorkDays"`
	OvertimeWeekCount      int                `json:"overtimeWeekCount"`
	RiskScore              int                `json:"riskScore"`
	RiskLevel              string             `json:"riskLevel"`
}

type ScheduleInsightMetricsResDto struct {
	RequiredHeadcount     int64                           `json:"requiredHeadcount"`
	AssignedHeadcount     int64                           `json:"assignedHeadcount"`
	UnfilledHeadcount     int64                           `json:"unfilledHeadcount"`
	CoverageRate          float64                         `json:"coverageRate"`
	OpenSwapRequestCount  int64                           `json:"openSwapRequestCount"`
	AverageHours          float64                         `json:"averageHours"`
	WorkloadSpreadHours   float64                         `json:"workloadSpreadHours"`
	EmployeesAtRisk       int                             `json:"employeesAtRisk"`
	AvailabilityConflicts int                             `json:"availabilityConflicts"`
	Employees             []ScheduleInsightEmployeeResDto `json:"employees"`
}

type ScheduleInsightAIUsageResDto struct {
	Used      int32     `json:"used"`
	Limit     int32     `json:"limit"`
	Remaining int32     `json:"remaining"`
	ResetAt   time.Time `json:"resetAt"`
}

type ScheduleInsightResDto struct {
	CompanyId   uuid.UUID                    `json:"companyId"`
	CompanyName string                       `json:"companyName"`
	StartAt     time.Time                    `json:"startAt"`
	EndAt       time.Time                    `json:"endAt"`
	Timezone    string                       `json:"timezone"`
	Locale      string                       `json:"locale"`
	Model       string                       `json:"model"`
	Workflow    []string                     `json:"workflow"`
	Metrics     ScheduleInsightMetricsResDto `json:"metrics"`
	AIUsage     ScheduleInsightAIUsageResDto `json:"aiUsage"`
	Summary     string                       `json:"summary"`
	GeneratedAt time.Time                    `json:"generatedAt"`
}
