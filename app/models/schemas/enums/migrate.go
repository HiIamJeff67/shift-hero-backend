package enums

// place the enums here to migrate
var MigratingEnums = map[string][]string{
	new(AccessControlPermission).Name():   AllAccessControlPermissionStrings,
	new(BillingIntervalUnit).Name():       AllBillingIntervalUnitStrings,
	new(BillingPlanName).Name():           AllBillingPlanNameStrings,
	new(BillingPlanStatus).Name():         AllBillingPlanStatusStrings,
	new(CountryCode).Name():               AllCountryCodeStrings,
	new(Country).Name():                   AllCountryStrings,
	new(Language).Name():                  AllLanguageStrings,
	new(SupportedCurrencyCode).Name():     AllSupportedCurrencyCodeStrings,
	new(UserGender).Name():                AllUserGenderStrings,
	new(UserPlan).Name():                  AllUserPlanStrings,
	new(UserRole).Name():                  AllUserRoleStrings,
	new(UserStatus).Name():                AllUserStatusStrings,
	new(UsersToBillingPlansStatus).Name(): AllUsersToBillingPlansStatusStrings,
}
