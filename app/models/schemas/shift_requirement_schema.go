package schemas

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type ShiftRequirement struct {
	Id            uuid.UUID          `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	CompanyId     uuid.UUID          `json:"companyId" gorm:"column:company_id; type:uuid; not null; index;"`
	EmployeeRole  enums.EmployeeRole `json:"employeeRole" gorm:"column:employee_role; type:\"EmployeeRole\"; not null; default:'Staff';"`
	StartAt       time.Time          `json:"startAt" gorm:"column:start_at; type:timestamptz; not null;"`
	EndAt         time.Time          `json:"endAt" gorm:"column:end_at; type:timestamptz; not null;"`
	RequiredCount int32              `json:"requiredCount" gorm:"column:required_count; type:integer; not null; default:1;"`
	Note          string             `json:"note" gorm:"column:note; size:1024; not null; default:'';"`
	UpdatedAt     time.Time          `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt     time.Time          `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Company Company `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (ShiftRequirement) TableName() string {
	return types.TableName_ShiftRequirementsTable.String()
}

type ShiftRequirementRelation types.RelationName

const (
	ShiftRequirementRelation_Company ShiftRequirementRelation = "Company"
)
