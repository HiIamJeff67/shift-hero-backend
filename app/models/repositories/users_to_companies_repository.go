package repositories

import (
	"time"

	"gorm.io/gorm/clause"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/shift-hero-backend/app/options"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

type CompanyMemberWithUser struct {
	UserId       uuid.UUID          `gorm:"column:user_id"`
	Name         string             `gorm:"column:name"`
	DisplayName  string             `gorm:"column:display_name"`
	Email        string             `gorm:"column:email"`
	EmployeeRole enums.EmployeeRole `gorm:"column:employee_role"`
}

type UsersToCompaniesRepositoryInterface interface {
	GetOneByCompanyIdAndUserId(companyId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UsersToCompanies, *exceptions.Exception)
	CreateOne(membership *schemas.UsersToCompanies, opts ...options.RepositoryOptions) *exceptions.Exception
	GetMembersByCompanyId(companyId uuid.UUID, opts ...options.RepositoryOptions) ([]CompanyMemberWithUser, *exceptions.Exception)
	CountManagersByCompanyId(companyId uuid.UUID, opts ...options.RepositoryOptions) (int64, *exceptions.Exception)
	UpdateEmployeeRole(companyId uuid.UUID, userId uuid.UUID, role enums.EmployeeRole, opts ...options.RepositoryOptions) (*time.Time, *exceptions.Exception)
	DeleteOneByCompanyIdAndUserId(companyId uuid.UUID, userId uuid.UUID, opts ...options.RepositoryOptions) *exceptions.Exception
}

type UsersToCompaniesRepository struct{}

func NewUsersToCompaniesRepository() UsersToCompaniesRepositoryInterface {
	return &UsersToCompaniesRepository{}
}

func (r *UsersToCompaniesRepository) GetOneByCompanyIdAndUserId(
	companyId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToCompanies, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	membership := schemas.UsersToCompanies{}
	result := parsedOptions.DB.Model(&schemas.UsersToCompanies{}).
		Where("company_id = ? AND user_id = ?", companyId, userId).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&membership)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Company.NotFound("Company membership not found").WithOrigin(result.Error)},
		{First: membership.CompanyId == uuid.Nil || membership.UserId == uuid.Nil, Second: exceptions.Company.NotFound("Company membership not found")},
	}); exception != nil {
		return nil, exception
	}

	return &membership, nil
}

func (r *UsersToCompaniesRepository) CreateOne(
	membership *schemas.UsersToCompanies,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if err := parsedOptions.DB.Model(&schemas.UsersToCompanies{}).Create(membership).Error; err != nil {
		return exceptions.Company.FailedToCreate().WithOrigin(err)
	}

	return nil
}

func (r *UsersToCompaniesRepository) GetMembersByCompanyId(
	companyId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]CompanyMemberWithUser, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	rows := []CompanyMemberWithUser{}
	result := parsedOptions.DB.Model(&schemas.UsersToCompanies{}).
		Select("\"UsersToCompaniesTable\".user_id, u.name, u.display_name, u.email, \"UsersToCompaniesTable\".employee_role").
		Joins("JOIN \"UserTable\" u ON u.id = \"UsersToCompaniesTable\".user_id").
		Where("\"UsersToCompaniesTable\".company_id = ?", companyId).
		Order("u.created_at ASC").
		Find(&rows)
	if result.Error != nil {
		return nil, exceptions.Company.NotFound().WithOrigin(result.Error)
	}

	return rows, nil
}

func (r *UsersToCompaniesRepository) CountManagersByCompanyId(
	companyId uuid.UUID,
	opts ...options.RepositoryOptions,
) (int64, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	count := int64(0)
	if err := parsedOptions.DB.Model(&schemas.UsersToCompanies{}).
		Where("company_id = ? AND employee_role = ?", companyId, enums.EmployeeRole_Manager).
		Count(&count).Error; err != nil {
		return 0, exceptions.Company.NotFound().WithOrigin(err)
	}

	return count, nil
}

func (r *UsersToCompaniesRepository) UpdateEmployeeRole(
	companyId uuid.UUID,
	userId uuid.UUID,
	role enums.EmployeeRole,
	opts ...options.RepositoryOptions,
) (*time.Time, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.UsersToCompanies{}).
		Where("company_id = ? AND user_id = ?", companyId, userId).
		Updates(map[string]any{"employee_role": role})
	if result.Error != nil {
		return nil, exceptions.Company.FailedToUpdate().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.Company.NotFound("Company member not found")
	}

	updated := schemas.UsersToCompanies{}
	if err := parsedOptions.DB.Model(&schemas.UsersToCompanies{}).
		Where("company_id = ? AND user_id = ?", companyId, userId).
		First(&updated).Error; err != nil {
		return nil, exceptions.Company.NotFound("Company member not found").WithOrigin(err)
	}

	return &updated.UpdatedAt, nil
}

func (r *UsersToCompaniesRepository) DeleteOneByCompanyIdAndUserId(
	companyId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.UsersToCompanies{}).
		Where("company_id = ? AND user_id = ?", companyId, userId).
		Delete(&schemas.UsersToCompanies{})
	if result.Error != nil {
		return exceptions.Company.FailedToDelete().WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return exceptions.Company.NotFound("Company member not found")
	}

	return nil
}
