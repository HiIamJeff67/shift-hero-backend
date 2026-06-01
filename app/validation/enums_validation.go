package validation

import (
	"github.com/go-playground/validator/v10" // make sure we use the version 10

	enums "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/enums"
	util "github.com/HiIamJeff67/shift-hero-backend/app/util"
)

func RegisterEnumsValidation(validate *validator.Validate) {
	validate.RegisterValidation("isaccesscontrolpermission", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllAccessControlPermissionStrings)
	})
	validate.RegisterValidation("isbillingintervalunit", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllBillingIntervalUnitStrings)
	})
	validate.RegisterValidation("isbillingplanname", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllBillingPlanNameStrings)
	})
	validate.RegisterValidation("isbillingplanstatus", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllBillingPlanStatusStrings)
	})
	validate.RegisterValidation("iscountrycode", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllCountryCodeStrings)
	})
	validate.RegisterValidation("iscountry", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllCountryStrings)
	})
	validate.RegisterValidation("islanguage", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllLanguageStrings)
	})
	validate.RegisterValidation("issupportedcurrencycode", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllSupportedCurrencyCodeStrings)
	})
	validate.RegisterValidation("isgender", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserGenderStrings)
	})
	validate.RegisterValidation("isplan", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserPlanStrings)
	})
	validate.RegisterValidation("isrole", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserRoleStrings)
	})
	validate.RegisterValidation("isstatus", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUserStatusStrings)
	})
	validate.RegisterValidation("isuserstobillingplansstatus", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		return util.IsStringIn(val, enums.AllUsersToBillingPlansStatusStrings)
	})
}
