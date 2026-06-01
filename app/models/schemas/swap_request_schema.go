package schemas

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type SwapRequest struct {
	Id                uuid.UUID               `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	CompanyId         uuid.UUID               `json:"companyId" gorm:"column:company_id; type:uuid; not null; index;"`
	ShiftAssignmentId uuid.UUID               `json:"shiftAssignmentId" gorm:"column:shift_assignment_id; type:uuid; not null; index;"`
	RequesterUserId   uuid.UUID               `json:"requesterUserId" gorm:"column:requester_user_id; type:uuid; not null; index;"`
	ClaimedByUserId   *uuid.UUID              `json:"claimedByUserId" gorm:"column:claimed_by_user_id; type:uuid; index;"`
	Status            enums.SwapRequestStatus `json:"status" gorm:"column:status; type:\"SwapRequestStatus\"; not null; default:'Open';"`
	Reason            string                  `json:"reason" gorm:"column:reason; size:1024; not null; default:'';"`
	UpdatedAt         time.Time               `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt         time.Time               `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Company         Company         `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ShiftAssignment ShiftAssignment `json:"shiftAssignment" gorm:"foreignKey:ShiftAssignmentId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RequesterUser   User            `json:"requesterUser" gorm:"foreignKey:RequesterUserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ClaimedByUser   User            `json:"claimedByUser" gorm:"foreignKey:ClaimedByUserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:SET NULL;"`
}

func (SwapRequest) TableName() string {
	return types.TableName_SwapRequestsTable.String()
}

type SwapRequestRelation types.RelationName

const (
	SwapRequestRelation_Company         SwapRequestRelation = "Company"
	SwapRequestRelation_ShiftAssignment SwapRequestRelation = "ShiftAssignment"
	SwapRequestRelation_RequesterUser   SwapRequestRelation = "RequesterUser"
	SwapRequestRelation_ClaimedByUser   SwapRequestRelation = "ClaimedByUser"
)
