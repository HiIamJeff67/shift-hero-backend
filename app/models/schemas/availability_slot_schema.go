package schemas

import (
	"time"

	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type AvailabilitySlot struct {
	Id          uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	CompanyId   uuid.UUID `json:"companyId" gorm:"column:company_id; type:uuid; not null; index;"`
	UserId      uuid.UUID `json:"userId" gorm:"column:user_id; type:uuid; not null; index;"`
	StartAt     time.Time `json:"startAt" gorm:"column:start_at; type:timestamptz; not null; index;"`
	EndAt       time.Time `json:"endAt" gorm:"column:end_at; type:timestamptz; not null; index;"`
	IsAvailable bool      `json:"isAvailable" gorm:"column:is_available; not null; default:true;"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Company Company `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	User    User    `json:"user" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (AvailabilitySlot) TableName() string {
	return types.TableName_AvailabilitySlotsTable.String()
}

type AvailabilitySlotRelation types.RelationName

const (
	AvailabilitySlotRelation_Company AvailabilitySlotRelation = "Company"
	AvailabilitySlotRelation_User    AvailabilitySlotRelation = "User"
)
