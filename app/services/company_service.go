package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/shift-hero-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	models "github.com/HiIamJeff67/shift-hero-backend/app/models"
	repositories "github.com/HiIamJeff67/shift-hero-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/shift-hero-backend/app/options"
	validation "github.com/HiIamJeff67/shift-hero-backend/app/validation"
)

type CompanyServiceInterface interface {
	CreateCompany(ctx context.Context, reqDto *dtos.CreateCompanyReqDto) (*dtos.CompanyResDto, *exceptions.Exception)
	GetMyCompanies(ctx context.Context, reqDto *dtos.GetMyCompaniesReqDto) ([]dtos.CompanyResDto, *exceptions.Exception)
	GetCompany(ctx context.Context, reqDto *dtos.GetCompanyReqDto) (*dtos.CompanyResDto, *exceptions.Exception)
	UpdateCompany(ctx context.Context, reqDto *dtos.UpdateCompanyReqDto) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception)
	GetCompanyMembers(ctx context.Context, reqDto *dtos.GetCompanyMembersReqDto) ([]dtos.CompanyMemberResDto, *exceptions.Exception)
	AddCompanyMember(ctx context.Context, reqDto *dtos.AddCompanyMemberReqDto) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception)
	UpdateCompanyMember(ctx context.Context, reqDto *dtos.UpdateCompanyMemberReqDto) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception)
	DeleteCompanyMember(ctx context.Context, reqDto *dtos.DeleteCompanyMemberReqDto) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception)
}

type CompanyService struct {
	db                         *gorm.DB
	companyRepository          repositories.CompanyRepositoryInterface
	usersToCompaniesRepository repositories.UsersToCompaniesRepositoryInterface
	userRepository             repositories.UserRepositoryInterface
}

func NewCompanyService(
	db *gorm.DB,
	companyRepository repositories.CompanyRepositoryInterface,
	usersToCompaniesRepository repositories.UsersToCompaniesRepositoryInterface,
	userRepository repositories.UserRepositoryInterface,
) CompanyServiceInterface {
	if db == nil {
		db = models.DB
	}
	return &CompanyService{
		db:                         db,
		companyRepository:          companyRepository,
		usersToCompaniesRepository: usersToCompaniesRepository,
		userRepository:             userRepository,
	}
}

func convertCompanyToRes(company schemas.Company) dtos.CompanyResDto {
	return dtos.CompanyResDto{
		Id:          company.Id,
		Name:        company.Name,
		Description: company.Description,
		Email:       company.Email,
		UpdatedAt:   company.UpdatedAt,
		CreatedAt:   company.CreatedAt,
	}
}

func (s *CompanyService) CreateCompany(
	ctx context.Context,
	reqDto *dtos.CreateCompanyReqDto,
) (*dtos.CompanyResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	tx := db.Begin()
	if tx.Error != nil {
		return nil, exceptions.Company.FailedToCommitTransaction().WithOrigin(tx.Error)
	}
	defer func() {
		if recover() != nil {
			tx.Rollback()
		}
	}()

	company := schemas.Company{
		Name:        reqDto.Body.Name,
		Description: reqDto.Body.Description,
		Email:       reqDto.Body.Email,
	}
	if exception := s.companyRepository.CreateOne(&company, options.WithDB(tx)); exception != nil {
		tx.Rollback()
		return nil, exception
	}

	member := schemas.UsersToCompanies{
		UserId:       reqDto.ContextFields.UserId,
		CompanyId:    company.Id,
		EmployeeRole: enums.EmployeeRole_Manager,
	}
	if exception := s.usersToCompaniesRepository.CreateOne(&member, options.WithDB(tx)); exception != nil {
		tx.Rollback()
		return nil, exceptions.Company.FailedToCreate("Failed to create owner membership").WithOrigin(exception.Origin)
	}

	setting := schemas.CompanySettings{CompanyId: company.Id}
	if err := tx.Model(&schemas.CompanySettings{}).Create(&setting).Error; err != nil {
		tx.Rollback()
		return nil, exceptions.Company.FailedToCreate("Failed to create default settings").WithOrigin(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.Company.FailedToCommitTransaction().WithOrigin(err)
	}

	res := convertCompanyToRes(company)
	return &res, nil
}

func (s *CompanyService) GetMyCompanies(
	ctx context.Context,
	reqDto *dtos.GetMyCompaniesReqDto,
) ([]dtos.CompanyResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	companies, exception := s.companyRepository.GetManyByUserId(
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	res := make([]dtos.CompanyResDto, len(companies))
	for i, company := range companies {
		res[i] = convertCompanyToRes(company)
	}

	return res, nil
}

func (s *CompanyService) GetCompany(
	ctx context.Context,
	reqDto *dtos.GetCompanyReqDto,
) (*dtos.CompanyResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.BadRequest("Invalid request payload").WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMemberByRepository(
		s.usersToCompaniesRepository,
		reqDto.Param.CompanyId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	company, exception := s.companyRepository.GetOneById(
		reqDto.Param.CompanyId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exceptions.Company.NotFound("Company not found for this companyId").WithOrigin(exception.Origin).WithDetails(map[string]any{
			"companyId": reqDto.Param.CompanyId.String(),
			"userId":    reqDto.ContextFields.UserId.String(),
		})
	}

	res := convertCompanyToRes(*company)
	return &res, nil
}

func (s *CompanyService) UpdateCompany(
	ctx context.Context,
	reqDto *dtos.UpdateCompanyReqDto,
) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManagerByRepository(
		s.usersToCompaniesRepository,
		reqDto.Body.CompanyId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	updates := map[string]any{}
	if reqDto.Body.Values.Name != nil {
		updates["name"] = *reqDto.Body.Values.Name
	}
	if reqDto.Body.Values.Description != nil {
		updates["description"] = *reqDto.Body.Values.Description
	}
	if reqDto.Body.Values.Email != nil {
		updates["email"] = *reqDto.Body.Values.Email
	}

	if len(updates) == 0 {
		return nil, exceptions.Company.NoChanges()
	}

	rowsAffected, exception := s.companyRepository.UpdateOneById(
		reqDto.Body.CompanyId,
		updates,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}
	if rowsAffected == 0 {
		return nil, exceptions.Company.NoChanges()
	}

	updated, exception := s.companyRepository.GetOneById(
		reqDto.Body.CompanyId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.MutationUpdatedAtResDto{UpdatedAt: updated.UpdatedAt}, nil
}

func (s *CompanyService) GetCompanyMembers(
	ctx context.Context,
	reqDto *dtos.GetCompanyMembersReqDto,
) ([]dtos.CompanyMemberResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyMemberByRepository(
		s.usersToCompaniesRepository,
		reqDto.Param.CompanyId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	rows, exception := s.usersToCompaniesRepository.GetMembersByCompanyId(
		reqDto.Param.CompanyId,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	res := make([]dtos.CompanyMemberResDto, len(rows))
	for i, row := range rows {
		res[i] = dtos.CompanyMemberResDto{
			UserId:       row.UserId,
			Name:         row.Name,
			DisplayName:  row.DisplayName,
			Email:        row.Email,
			EmployeeRole: row.EmployeeRole,
		}
	}
	return res, nil
}

func (s *CompanyService) AddCompanyMember(
	ctx context.Context,
	reqDto *dtos.AddCompanyMemberReqDto,
) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManagerByRepository(
		s.usersToCompaniesRepository,
		reqDto.Body.CompanyId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	if _, exception := s.userRepository.GetOneById(
		reqDto.Body.UserId,
		nil,
		options.WithDB(db),
	); exception != nil {
		return nil, exceptions.User.NotFound().WithOrigin(exception.Origin)
	}

	if _, exception := s.usersToCompaniesRepository.GetOneByCompanyIdAndUserId(
		reqDto.Body.CompanyId,
		reqDto.Body.UserId,
		options.WithDB(db),
	); exception == nil {
		return nil, exceptions.Company.DuplicateMember(reqDto.Body.CompanyId.String(), reqDto.Body.UserId.String())
	} else if exception.Reason != "NotFound" {
		return nil, exception
	}

	member := schemas.UsersToCompanies{
		CompanyId:    reqDto.Body.CompanyId,
		UserId:       reqDto.Body.UserId,
		EmployeeRole: reqDto.Body.EmployeeRole,
	}
	if exception := s.usersToCompaniesRepository.CreateOne(
		&member,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	return &dtos.MutationUpdatedAtResDto{UpdatedAt: member.UpdatedAt}, nil
}

func (s *CompanyService) UpdateCompanyMember(
	ctx context.Context,
	reqDto *dtos.UpdateCompanyMemberReqDto,
) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManagerByRepository(
		s.usersToCompaniesRepository,
		reqDto.Body.CompanyId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	updatedAt, exception := s.usersToCompaniesRepository.UpdateEmployeeRole(
		reqDto.Body.CompanyId,
		reqDto.Body.UserId,
		reqDto.Body.EmployeeRole,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.MutationUpdatedAtResDto{UpdatedAt: *updatedAt}, nil
}

func (s *CompanyService) DeleteCompanyMember(
	ctx context.Context,
	reqDto *dtos.DeleteCompanyMemberReqDto,
) (*dtos.MutationUpdatedAtResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Company.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)
	if _, exception := requireCompanyManagerByRepository(
		s.usersToCompaniesRepository,
		reqDto.Body.CompanyId,
		reqDto.ContextFields.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	if reqDto.Body.UserId == reqDto.ContextFields.UserId {
		managerCount, exception := s.usersToCompaniesRepository.CountManagersByCompanyId(
			reqDto.Body.CompanyId,
			options.WithDB(db),
		)
		if exception != nil {
			return nil, exceptions.Company.FailedToDelete().WithOrigin(exception.Origin)
		}
		if managerCount <= 1 {
			return nil, exceptions.Company.Forbidden("Cannot remove the last manager from company")
		}
	}

	if exception := s.usersToCompaniesRepository.DeleteOneByCompanyIdAndUserId(
		reqDto.Body.CompanyId,
		reqDto.Body.UserId,
		options.WithDB(db),
	); exception != nil {
		return nil, exception
	}

	return &dtos.MutationUpdatedAtResDto{UpdatedAt: truncateToMinute(time.Now().UTC())}, nil
}
