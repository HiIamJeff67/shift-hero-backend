package repositories

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "github.com/HiIamJeff67/shift-hero-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/shift-hero-backend/app/models/inputs"
	schemas "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas"
	useraccountsqls "github.com/HiIamJeff67/shift-hero-backend/app/models/sqls/user_account"
	options "github.com/HiIamJeff67/shift-hero-backend/app/options"
	util "github.com/HiIamJeff67/shift-hero-backend/app/util"
	types "github.com/HiIamJeff67/shift-hero-backend/shared/types"
)

type UserAccountRepositoryInterface interface {
	GetOneByUserId(userId uuid.UUID, opts ...options.RepositoryOptions) (*schemas.UserAccount, *exceptions.Exception)
	GetAIUsageQuotaByUserIdForUpdate(userId uuid.UUID, opts ...options.RepositoryOptions) (*UserAIUsageQuota, *exceptions.Exception)
	CreateOneByUserId(userId uuid.UUID, input inputs.CreateUserAccountInput, opts ...options.RepositoryOptions) (*uuid.UUID, *exceptions.Exception)
	UpdateOneByUserId(userId uuid.UUID, input inputs.PartialUpdateUserAccountInput, opts ...options.RepositoryOptions) (*schemas.UserAccount, *exceptions.Exception)
	UpdateAIUsageByUserId(userId uuid.UUID, input inputs.UpdateUserAIUsageInput, opts ...options.RepositoryOptions) *exceptions.Exception
	ReleaseAIUsageReservationByUserId(userId uuid.UUID, periodStart time.Time, opts ...options.RepositoryOptions) *exceptions.Exception
}

type UserAccountRepository struct{}

type UserAIUsageQuota struct {
	UserId            uuid.UUID `gorm:"column:user_id"`
	MonthlyUsageCount int32     `gorm:"column:monthly_usage_count"`
	PeriodStart       time.Time `gorm:"column:period_start"`
	MonthlyLimit      int32     `gorm:"column:monthly_limit"`
}

func NewUserAccountRepository() UserAccountRepositoryInterface {
	return &UserAccountRepository{}
}

func (r *UserAccountRepository) GetOneByUserId(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*schemas.UserAccount, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var userAccount schemas.UserAccount
	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Where("user_id = ?", userId).
		Clauses(clause.Locking{Strength: "SHARE"}).
		First(&userAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.NotFound().WithOrigin(err)
	}

	return &userAccount, nil
}

func (r *UserAccountRepository) GetAIUsageQuotaByUserIdForUpdate(
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (*UserAIUsageQuota, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var quota UserAIUsageQuota
	result := parsedOptions.DB.Raw(
		useraccountsqls.GetAIUsageQuotaByUserIdForUpdateSQL,
		sql.Named("user_id", userId),
	).Scan(&quota)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToGetAIUsageQuota().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserAccount.NotFound()
	}
	return &quota, nil
}

func (r *UserAccountRepository) CreateOneByUserId(
	userId uuid.UUID,
	input inputs.CreateUserAccountInput,
	opts ...options.RepositoryOptions,
) (*uuid.UUID, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	var newUserAccount schemas.UserAccount
	newUserAccount.UserId = userId

	if err := copier.Copy(&newUserAccount, &input); err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithOrigin(err)
	}

	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&newUserAccount)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToCreate().WithOrigin(err)
	}

	return &newUserAccount.Id, nil
}

func (r *UserAccountRepository) UpdateOneByUserId(
	userId uuid.UUID,
	input inputs.PartialUpdateUserAccountInput,
	opts ...options.RepositoryOptions,
) (*schemas.UserAccount, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)

	existingUserAccount, exception := r.GetOneByUserId(
		userId,
		opts...,
	)
	if exception = exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
		{First: existingUserAccount == nil, Second: exceptions.UserAccount.NotFound()},
	}); exception != nil {
		return nil, exception
	}

	updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingUserAccount)
	if err != nil {
		return nil, exceptions.Util.FailedToPreprocessPartialUpdate(input.Values, input.SetNull, *existingUserAccount)
	}

	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Where("user_id = ?", userId).
		Select("*").
		Updates(&updates)
	if err := result.Error; err != nil {
		return nil, exceptions.UserAccount.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.UserAccount.NoChanges()
	}

	return &updates, nil
}

func (r *UserAccountRepository) UpdateAIUsageByUserId(
	userId uuid.UUID,
	input inputs.UpdateUserAIUsageInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Where("user_id = ?", userId).
		Updates(map[string]any{
			"ai_monthly_usage_count": input.AIMonthlyUsageCount,
			"ai_usage_period_start":  input.AIUsagePeriodStart,
		})
	if err := result.Error; err != nil {
		return exceptions.UserAccount.FailedToUpdate().WithOrigin(err)
	}
	if result.RowsAffected == 0 {
		return exceptions.UserAccount.NotFound()
	}
	return nil
}

func (r *UserAccountRepository) ReleaseAIUsageReservationByUserId(
	userId uuid.UUID,
	periodStart time.Time,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	result := parsedOptions.DB.Model(&schemas.UserAccount{}).
		Where(
			"user_id = ? AND ai_usage_period_start = ? AND ai_monthly_usage_count > 0",
			userId,
			periodStart,
		).
		UpdateColumn("ai_monthly_usage_count", gorm.Expr("ai_monthly_usage_count - 1"))
	if err := result.Error; err != nil {
		return exceptions.UserAccount.FailedToUpdate().WithOrigin(err)
	}
	return nil
}

// We do not allow to just delete the userAccount,
// instead, the userAccount is only deleted by deleting the user
// func DeleteUserAccount(userId uuid.UUID) (deletedUserAccount User, err error) {}
