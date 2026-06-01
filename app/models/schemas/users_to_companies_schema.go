package schemas

import (
	"time"

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type UsersToCompanies struct {
	UserId       uuid.UUID          `json:"userId" gorm:"column:user_id; type:uuid; primaryKey;"`
	CompanyId    uuid.UUID          `json:"companyId" gorm:"column:company_id; type:uuid; primaryKey;"`
	EmployeeRole enums.EmployeeRole `json:"employeeRole" gorm:"column:employee_role; type:\"EmployeeRole\"; not null; default:'Staff';"`
	UpdatedAt    time.Time          `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt    time.Time          `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	User    User    `json:"user" gorm:"foreignKey:UserId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Company Company `json:"company" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// UsersToCompanies Table Name
func (UsersToCompanies) TableName() string {
	return types.TableName_UsersToCompaniesTable.String()
}

// UsersToCompanies Relations
type UsersToCompaniesRelation types.RelationName

const (
	UsersToCompaniesRelation_User    UsersToCompaniesRelation = "User"
	UsersToCompaniesRelation_Company UsersToCompaniesRelation = "Company"
)
