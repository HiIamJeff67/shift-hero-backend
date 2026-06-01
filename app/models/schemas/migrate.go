package schemas

// place the tables here to migrate
var MigratingTables = []any{
	// public tables
	&User{},
	&UserInfo{},
	&UserAccount{},
	&UserSetting{},

	&UsersToBillingPlans{},
	&Company{},
	&UsersToCompanies{},
	&CompanySettings{},
	&ShiftRequirement{},
	&AvailabilitySlot{},
	&ShiftAssignment{},
	&SwapRequest{},

	// private tables
	&PlanLimitation{},
	&BillingPlan{},
}
