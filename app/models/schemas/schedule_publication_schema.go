package schemas

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type SchedulePublication struct {
	Id                uuid.UUID                       `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	CompanyId         uuid.UUID                       `json:"companyId" gorm:"column:company_id; type:uuid; not null; index;"`
	WeekStart         time.Time                       `json:"weekStart" gorm:"column:week_start; type:date; not null; index;"`
	Timezone          string                          `json:"timezone" gorm:"column:timezone; size:64; not null; default:'Asia/Taipei';"`
	Status            enums.SchedulePublicationStatus `json:"status" gorm:"column:status; type:\"SchedulePublicationStatus\"; not null; default:'Draft';"`
	PublishedByUserId *uuid.UUID                      `json:"publishedByUserId" gorm:"column:published_by_user_id; type:uuid; index;"`
	PublishedAt       *time.Time                      `json:"publishedAt" gorm:"column:published_at; type:timestamptz;"`
	UpdatedAt         time.Time                       `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt         time.Time                       `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Company         Company `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	PublishedByUser User    `json:"publishedByUser" gorm:"foreignKey:PublishedByUserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:SET NULL;"`
}

func (SchedulePublication) TableName() string {
	return types.TableName_SchedulePublicationsTable.String()
}

type SchedulePublicationRelation types.RelationName

const (
	SchedulePublicationRelation_Company         SchedulePublicationRelation = "Company"
	SchedulePublicationRelation_PublishedByUser SchedulePublicationRelation = "PublishedByUser"
)
