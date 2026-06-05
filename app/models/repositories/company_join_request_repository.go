package repositories

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/shift-hero-backend/app/options"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

type CompanyJoinRequestWithDetails struct {
	Id               uuid.UUID                      `gorm:"column:id"`
	CompanyId        uuid.UUID                      `gorm:"column:company_id"`
	CompanyName      string                         `gorm:"column:company_name"`
	RequesterUserId  uuid.UUID                      `gorm:"column:requester_user_id"`
	RequesterName    string                         `gorm:"column:requester_name"`
	RequesterEmail   string                         `gorm:"column:requester_email"`
	RequestedRole    enums.EmployeeRole             `gorm:"column:requested_role"`
	Note             string                         `gorm:"column:note"`
	Status           enums.CompanyJoinRequestStatus `gorm:"column:status"`
	ReviewedByUserId *uuid.UUID                     `gorm:"column:reviewed_by_user_id"`
	ReviewedAt       *time.Time                     `gorm:"column:reviewed_at"`
	CreatedAt        time.Time                      `gorm:"column:created_at"`
	UpdatedAt        time.Time                      `gorm:"column:updated_at"`
}

type CompanyJoinRequestRepositoryInterface interface {
	CreateOne(joinRequest *schemas.CompanyJoinRequest, opts ...options.RepositoryOptions) *exceptions.Exception
	GetOneByIdAndCompanyId(id uuid.UUID, companyId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.CompanyJoinRequest, *exceptions.Exception)
	GetOneWithDetailsById(id uuid.UUID, opts ...options.RepositoryOptions) (*CompanyJoinRequestWithDetails, *exceptions.Exception)
	GetManyWithDetailsByCompanyId(companyId uuid.UUID, status *enums.CompanyJoinRequestStatus, opts ...options.RepositoryOptions) ([]CompanyJoinRequestWithDetails, *exceptions.Exception)
	GetManyWithDetailsByRequesterUserId(requesterUserId uuid.UUID, opts ...options.RepositoryOptions) ([]CompanyJoinRequestWithDetails, *exceptions.Exception)
	GetPendingByCompanyIdAndRequesterUserId(companyId uuid.UUID, requesterUserId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.CompanyJoinRequest, *exceptions.Exception)
	UpdateReviewState(id uuid.UUID, status enums.CompanyJoinRequestStatus, reviewedByUserId uuid.UUID, reviewedAt time.Time, opts ...options.RepositoryOptions) *exceptions.Exception
}

type CompanyJoinRequestRepository struct{}

func NewCompanyJoinRequestRepository() CompanyJoinRequestRepositoryInterface {
	return &CompanyJoinRequestRepository{}
}

func (r *CompanyJoinRequestRepository) CreateOne(
	joinRequest *schemas.CompanyJoinRequest,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	if err := parsedOptions.DB.Model(&schemas.CompanyJoinRequest{}).Create(joinRequest).Error; err != nil {
		if strings.Contains(err.Error(), "company_join_requests_idx_company_id_requester_user_id_pending") {
			return exceptions.Company.DuplicateJoinRequest(joinRequest.CompanyId.String(), joinRequest.RequesterUserId.String())
		}
		return exceptions.Company.FailedToCreate("Failed to create company join request").WithOrigin(err)
	}

	return nil
}

func (r *CompanyJoinRequestRepository) GetOneByIdAndCompanyId(
	id uuid.UUID,
	companyId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.CompanyJoinRequest, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	joinRequest := schemas.CompanyJoinRequest{}
	result := parsedOptions.DB.Model(&schemas.CompanyJoinRequest{}).
		Where("id = ? AND company_id = ?", id, companyId).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&joinRequest)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Company.NotFound("Company join request not found").WithOrigin(result.Error)},
		{First: joinRequest.Id == uuid.Nil, Second: exceptions.Company.NotFound("Company join request not found")},
	}); exception != nil {
		return nil, exception
	}

	return &joinRequest, nil
}

func (r *CompanyJoinRequestRepository) GetOneWithDetailsById(
	id uuid.UUID,
	opts ...options.RepositoryOptions,
) (*CompanyJoinRequestWithDetails, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	row := CompanyJoinRequestWithDetails{}
	result := buildCompanyJoinRequestDetailsQuery(parsedOptions.DB).
		Where("cjr.id = ?", id).
		First(&row)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Company.NotFound("Company join request not found").WithOrigin(result.Error)},
		{First: row.Id == uuid.Nil, Second: exceptions.Company.NotFound("Company join request not found")},
	}); exception != nil {
		return nil, exception
	}

	return &row, nil
}

func (r *CompanyJoinRequestRepository) GetManyWithDetailsByCompanyId(
	companyId uuid.UUID,
	status *enums.CompanyJoinRequestStatus,
	opts ...options.RepositoryOptions,
) ([]CompanyJoinRequestWithDetails, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	rows := []CompanyJoinRequestWithDetails{}
	query := buildCompanyJoinRequestDetailsQuery(parsedOptions.DB).
		Where("cjr.company_id = ?", companyId)
	if status != nil {
		query = query.Where("cjr.status = ?", *status)
	}
	if err := query.Order("cjr.created_at DESC").Find(&rows).Error; err != nil {
		return nil, exceptions.Company.NotFound("Company join requests not found").WithOrigin(err)
	}

	return rows, nil
}

func (r *CompanyJoinRequestRepository) GetManyWithDetailsByRequesterUserId(
	requesterUserId uuid.UUID,
	opts ...options.RepositoryOptions,
) ([]CompanyJoinRequestWithDetails, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	rows := []CompanyJoinRequestWithDetails{}
	if err := buildCompanyJoinRequestDetailsQuery(parsedOptions.DB).
		Where("cjr.requester_user_id = ?", requesterUserId).
		Order("cjr.created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, exceptions.Company.NotFound("Company join requests not found").WithOrigin(err)
	}

	return rows, nil
}

func (r *CompanyJoinRequestRepository) GetPendingByCompanyIdAndRequesterUserId(
	companyId uuid.UUID,
	requesterUserId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.CompanyJoinRequest, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	joinRequest := schemas.CompanyJoinRequest{}
	result := parsedOptions.DB.Model(&schemas.CompanyJoinRequest{}).
		Where("company_id = ? AND requester_user_id = ? AND status = ?", companyId, requesterUserId, enums.CompanyJoinRequestStatus_Pending).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&joinRequest)
	if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
		{First: result.Error != nil, Second: exceptions.Company.NotFound("Pending company join request not found").WithOrigin(result.Error)},
		{First: joinRequest.Id == uuid.Nil, Second: exceptions.Company.NotFound("Pending company join request not found")},
	}); exception != nil {
		return nil, exception
	}

	return &joinRequest, nil
}

func (r *CompanyJoinRequestRepository) UpdateReviewState(
	id uuid.UUID,
	status enums.CompanyJoinRequestStatus,
	reviewedByUserId uuid.UUID,
	reviewedAt time.Time,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	result := parsedOptions.DB.Model(&schemas.CompanyJoinRequest{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":              status,
			"reviewed_by_user_id": reviewedByUserId,
			"reviewed_at":         reviewedAt,
		})
	if result.Error != nil {
		return exceptions.Company.FailedToUpdate("Failed to update company join request").WithOrigin(result.Error)
	}
	if result.RowsAffected == 0 {
		return exceptions.Company.NotFound("Company join request not found")
	}

	return nil
}

func buildCompanyJoinRequestDetailsQuery(db *gorm.DB) *gorm.DB {
	return db.Table("\"CompanyJoinRequestsTable\" AS cjr").
		Select(`
			cjr.id,
			cjr.company_id,
			c.name AS company_name,
			cjr.requester_user_id,
			u.name AS requester_name,
			u.email AS requester_email,
			cjr.requested_role,
			cjr.note,
			cjr.status,
			cjr.reviewed_by_user_id,
			cjr.reviewed_at,
			cjr.created_at,
			cjr.updated_at
		`).
		Joins("JOIN \"CompanyTable\" c ON c.id = cjr.company_id").
		Joins("JOIN \"UserTable\" u ON u.id = cjr.requester_user_id")
}
