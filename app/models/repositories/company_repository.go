package repositories

import (
	"gorm.io/gorm/clause"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	options "github.com/HiIamJeff67/shift-hero-backend/app/options"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

type CompanyRepositoryInterface interface {
	GetOneById(id uuid.UUID, opts ...options.RepositoryOptions) (*schemas.Company, *exceptions.Exception)
	GetManyByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) ([]schemas.Company, *exceptions.Exception)
	CreateOne(company *schemas.Company, opts ...options.RepositoryOptions) *exceptions.Exception
	UpdateOneById(id uuid.UUID, updates map[string]any, opts ...options.RepositoryOptions) (int64, *exceptions.Exception)
}

type CompanyRepository struct{}

func NewCompanyRepository() CompanyRepositoryInterface {
	return &CompanyRepository{}
}

func (r *CompanyRepository) GetOneById(
	id uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.Company, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	company := schemas.Company{}
	result := parsedOptions.DB.Model(&schemas.Company{}).
		Where("id = ?", id).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&company)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Company.NotFound().WithOrigin(result.Error)},
		{First: company.Id == uuid.Nil, Second: exceptions.Company.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	return &company, nil
}

func (r *CompanyRepository) GetManyByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]schemas.Company, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	companies := []schemas.Company{}
	result := parsedOptions.DB.Model(&schemas.Company{}).
		Joins("JOIN \"UsersToCompaniesTable\" utc ON utc.company_id = \"CompanyTable\".id").
		Where("utc.user_id = ?", userId).
		Order("created_at DESC").
		Find(&companies)
	if result.Error != nil {
		return nil, exceptions.Company.NotFound().WithOrigin(result.Error)
	}

	return companies, nil
}

func (r *CompanyRepository) CreateOne(
	company *schemas.Company,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if err := parsedOptions.DB.Model(&schemas.Company{}).Create(company).Error; err != nil {
		return exceptions.Company.FailedToCreate().WithOrigin(err)
	}

	return nil
}

func (r *CompanyRepository) UpdateOneById(
	id uuid.UUID,
	updates map[string]any,
	opts ...options.RepositoryOptions,
) (int64, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.Company{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return 0, exceptions.Company.FailedToUpdate().WithOrigin(result.Error)
	}

	return result.RowsAffected, nil
}
