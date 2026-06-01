package schemas

import (
	"time"

	enums "github.com/your-org/go-start-monolithic-kit/app/models/schemas/enums"
	types "github.com/your-org/go-start-monolithic-kit/shared/types"
)

// This table is only mutatable by the admin, and accessable by both client user and admin.
// To declare the value or data of this table, you MUST use the seeding method under github.com/your-org/go-start-monolithic-kit/app/models/seeds/
type PlanLimitation struct {
	Key       enums.UserPlan `json:"key" gorm:"column:key; type:\"UserPlan\"; primaryKey;"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`
}

// PlanLimitation Table Name
func (PlanLimitation) TableName() string {
	return types.TableName_PlanLimitationTable.String()
}
