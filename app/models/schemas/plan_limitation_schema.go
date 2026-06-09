package schemas

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

// This table is only mutatable by the admin, and accessable by both client user and admin.
// To declare the value or data of this table, you MUST use the seeding method under github.com/HiIamJeff67/shift-hero-backend/app/models/seeds/
type PlanLimitation struct {
	Key                      enums.UserPlan `json:"key" gorm:"column:key; type:\"UserPlan\"; primaryKey;"`
	AIMonthlyGenerationLimit int32          `json:"aiMonthlyGenerationLimit" gorm:"column:ai_monthly_generation_limit; type:integer; not null; default:0;"`
	UpdatedAt                time.Time      `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt                time.Time      `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`
}

// PlanLimitation Table Name
func (PlanLimitation) TableName() string {
	return types.TableName_PlanLimitationTable.String()
}
