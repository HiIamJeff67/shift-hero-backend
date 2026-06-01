package services

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	repositories "github.com/HiIamJeff67/shift-hero-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/shift-hero-backend/app/options"
)

func truncateToMinute(t time.Time) time.Time {
	return t.Truncate(time.Minute)
}

func validateTimeRange(startAt time.Time, endAt time.Time) *exceptions.Exception {
	if !endAt.After(startAt) {
		return exceptions.Scheduling.InvalidTimeRange()
	}
	return nil
}

func findCompanyMembership(db *gorm.DB, companyId uuid.UUID, userId uuid.UUID) (*schemas.UsersToCompanies, *exceptions.Exception) {
	membership := schemas.UsersToCompanies{}
	result := db.Model(&schemas.UsersToCompanies{}).
		Where("company_id = ? AND user_id = ?", companyId, userId).
		First(&membership)
	if result.Error != nil {
		return nil, exceptions.Company.NotFound("Company membership not found").WithOrigin(result.Error)
	}
	return &membership, nil
}

func requireCompanyMember(db *gorm.DB, companyId uuid.UUID, userId uuid.UUID) (*schemas.UsersToCompanies, *exceptions.Exception) {
	membership, exception := findCompanyMembership(db, companyId, userId)
	if exception != nil {
		if exception.Origin != nil && exception.Reason != "NotFound" {
			return nil, exception.WithDetails(map[string]any{
				"companyId": companyId.String(),
				"userId":    userId.String(),
				"stage":     "findCompanyMembership",
			})
		}
		return nil, exceptions.Company.Forbidden("You are not a member of this company")
	}
	return membership, nil
}

func requireCompanyManager(db *gorm.DB, companyId uuid.UUID, userId uuid.UUID) (*schemas.UsersToCompanies, *exceptions.Exception) {
	membership, exception := requireCompanyMember(db, companyId, userId)
	if exception != nil {
		return nil, exception
	}
	if membership.EmployeeRole != enums.EmployeeRole_Manager {
		return nil, exceptions.Company.Forbidden("Manager role is required for this company operation")
	}
	return membership, nil
}

func findCompanyMembershipByRepository(
	repository repositories.UsersToCompaniesRepositoryInterface,
	companyId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToCompanies, *exceptions.Exception) {
	membership, exception := repository.GetOneByCompanyIdAndUserId(companyId, userId, opts...)
	if exception != nil {
		return nil, exception
	}

	return membership, nil
}

func requireCompanyMemberByRepository(
	repository repositories.UsersToCompaniesRepositoryInterface,
	companyId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToCompanies, *exceptions.Exception) {
	membership, exception := findCompanyMembershipByRepository(repository, companyId, userId, opts...)
	if exception != nil {
		if exception.Origin != nil && exception.Reason != "NotFound" {
			return nil, exception.WithDetails(map[string]any{
				"companyId": companyId.String(),
				"userId":    userId.String(),
				"stage":     "findCompanyMembershipByRepository",
			})
		}
		return nil, exceptions.Company.Forbidden("You are not a member of this company")
	}

	return membership, nil
}

func requireCompanyManagerByRepository(
	repository repositories.UsersToCompaniesRepositoryInterface,
	companyId uuid.UUID,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UsersToCompanies, *exceptions.Exception) {
	membership, exception := requireCompanyMemberByRepository(repository, companyId, userId, opts...)
	if exception != nil {
		return nil, exception
	}
	if membership.EmployeeRole != enums.EmployeeRole_Manager {
		return nil, exceptions.Company.Forbidden("Manager role is required for this company operation")
	}

	return membership, nil
}
