package schemas

// place the tables here to migrate
var MigratingTables = []any{
	// public tables
	&User{},
	&UserInfo{},
	&UserAccount{},
	&UserSetting{},

	&UsersToBillingPlans{},

	// private tables
	&PlanLimitation{},
	&BillingPlan{},
}
