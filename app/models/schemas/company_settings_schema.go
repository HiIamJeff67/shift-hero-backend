package schemas

import (
	"time"

	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type CompanySettings struct {
	CompanyId        uuid.UUID `json:"companyId" gorm:"column:company_id; type:uuid; primaryKey;"`
	AutoApproveSwaps bool      `json:"autoApproveSwaps" gorm:"column:auto_approve_swaps; not null; default:false;"`
	MaxWeeklyHours   int32     `json:"maxWeeklyHours" gorm:"column:max_weekly_hours; type:integer; not null; default:40;"`
	MinRestHours     int32     `json:"minRestHours" gorm:"column:min_rest_hours; type:integer; not null; default:8;"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt        time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Company Company `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (CompanySettings) TableName() string {
	return types.TableName_CompanySettingsTable.String()
}

type CompanySettingsRelation types.RelationName

const (
	CompanySettingsRelation_Company CompanySettingsRelation = "Company"
)
