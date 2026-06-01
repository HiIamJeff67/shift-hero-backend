package inputs

import (
	enums "github.com/your-org/go-start-monolithic-kit/app/models/schemas/enums"
)

type CreateUserSettingInput struct {
	Language           *enums.Language `json:"language" gorm:"column:language;"`
	GeneralSettingCode *int64          `json:"generalSettingCode" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int64          `json:"privacySettingCode" gorm:"column:privacy_setting_code;"`
}

type UpdateUserSettingInput struct {
	Language           *enums.Language `json:"language" gorm:"column:language;"`
	GeneralSettingCode *int64          `json:"generalSettingCode" gorm:"column:general_setting_code;"`
	PrivacySettingCode *int64          `json:"privacySettingCode" gorm:"column:privacy_setting_code;"`
}

type PartialUpdateUserSettingInput = PartialUpdateInput[UpdateUserSettingInput]
