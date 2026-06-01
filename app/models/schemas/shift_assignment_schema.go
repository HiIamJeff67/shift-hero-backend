package schemas

import (
	"time"

	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type ShiftAssignment struct {
	Id                 uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	CompanyId          uuid.UUID `json:"companyId" gorm:"column:company_id; type:uuid; not null; index;"`
	ShiftRequirementId uuid.UUID `json:"shiftRequirementId" gorm:"column:shift_requirement_id; type:uuid; not null; index;"`
	UserId             uuid.UUID `json:"userId" gorm:"column:user_id; type:uuid; not null; index;"`
	StartAt            time.Time `json:"startAt" gorm:"column:start_at; type:timestamptz; not null; index;"`
	EndAt              time.Time `json:"endAt" gorm:"column:end_at; type:timestamptz; not null; index;"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt          time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Company          Company          `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	User             User             `json:"user" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ShiftRequirement ShiftRequirement `json:"shiftRequirement" gorm:"foreignKey:ShiftRequirementId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (ShiftAssignment) TableName() string {
	return types.TableName_ShiftAssignmentsTable.String()
}

type ShiftAssignmentRelation types.RelationName

const (
	ShiftAssignmentRelation_Company          ShiftAssignmentRelation = "Company"
	ShiftAssignmentRelation_User             ShiftAssignmentRelation = "User"
	ShiftAssignmentRelation_ShiftRequirement ShiftAssignmentRelation = "ShiftRequirement"
)
