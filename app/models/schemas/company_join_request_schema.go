package schemas

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type CompanyJoinRequest struct {
	Id               uuid.UUID                      `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	CompanyId        uuid.UUID                      `json:"companyId" gorm:"column:company_id; type:uuid; not null; index;"`
	RequesterUserId  uuid.UUID                      `json:"requesterUserId" gorm:"column:requester_user_id; type:uuid; not null; index;"`
	RequestedRole    enums.EmployeeRole             `json:"requestedRole" gorm:"column:requested_role; type:\"EmployeeRole\"; not null; default:'Staff';"`
	Note             string                         `json:"note" gorm:"column:note; type:text; not null; default:'';"`
	Status           enums.CompanyJoinRequestStatus `json:"status" gorm:"column:status; type:\"CompanyJoinRequestStatus\"; not null; default:'Pending'; index;"`
	ReviewedByUserId *uuid.UUID                     `json:"reviewedByUserId" gorm:"column:reviewed_by_user_id; type:uuid; index;"`
	ReviewedAt       *time.Time                     `json:"reviewedAt" gorm:"column:reviewed_at; type:timestamptz;"`
	UpdatedAt        time.Time                      `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt        time.Time                      `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	Company        Company `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	RequesterUser  User    `json:"requesterUser" gorm:"foreignKey:RequesterUserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ReviewedByUser User    `json:"reviewedByUser" gorm:"foreignKey:ReviewedByUserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:SET NULL;"`
}

func (CompanyJoinRequest) TableName() string {
	return types.TableName_CompanyJoinRequestsTable.String()
}

type CompanyJoinRequestRelation types.RelationName

const (
	CompanyJoinRequestRelation_Company        CompanyJoinRequestRelation = "Company"
	CompanyJoinRequestRelation_RequesterUser  CompanyJoinRequestRelation = "RequesterUser"
	CompanyJoinRequestRelation_ReviewedByUser CompanyJoinRequestRelation = "ReviewedByUser"
)
