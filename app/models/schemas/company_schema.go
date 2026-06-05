package schemas

import (
	"time"

	"github.com/HiIamJeff67/shift-hero-backend/shared/types"
	"github.com/google/uuid"
)

type Company struct {
	Id          uuid.UUID `json:"id" gorm:"column:id; type:uuid; primaryKey; default:gen_random_uuid();"`
	Name        string    `json:"name" gorm:"column:name; size:128; unique; not null;"`
	Description string    `json:"description" gorm:"column:description; size:1024; not null; default:'';"`
	Email       string    `json:"email" gorm:"column:email; size:128; unique; not null;"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at; type:timestamptz; not null; autoUpdateTime:true;"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at; type:timestamptz; not null; autoCreateTime:true;"`

	// relations
	UsersToCompanies  []UsersToCompanies   `json:"usersToCompanies" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CompanySettings   *CompanySettings     `json:"companySettings" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ShiftRequirements []ShiftRequirement   `json:"shiftRequirements" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AvailabilitySlots []AvailabilitySlot   `json:"availabilitySlots" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ShiftAssignments  []ShiftAssignment    `json:"shiftAssignments" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	SwapRequests      []SwapRequest        `json:"swapRequests" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
	JoinRequests      []CompanyJoinRequest `json:"joinRequests" gorm:"foreignKey:CompanyId; references:Id; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

// Company Table Name
func (Company) TableName() string {
	return types.TableName_CompanyTable.String()
}

// Company Relations
type CompanyRelation types.RelationName

const (
	CompanyRelation_UsersToCompanies  CompanyRelation = "UsersToCompanies"
	CompanyRelation_CompanySettings   CompanyRelation = "CompanySettings"
	CompanyRelation_ShiftRequirements CompanyRelation = "ShiftRequirements"
	CompanyRelation_AvailabilitySlots CompanyRelation = "AvailabilitySlots"
	CompanyRelation_ShiftAssignments  CompanyRelation = "ShiftAssignments"
	CompanyRelation_SwapRequests      CompanyRelation = "SwapRequests"
	CompanyRelation_JoinRequests      CompanyRelation = "JoinRequests"
)
