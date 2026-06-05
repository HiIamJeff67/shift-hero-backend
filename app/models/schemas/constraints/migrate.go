package constraints

import (
	companyjoinrequestsconstraints "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/constraints/company_join_requests_constraints"
	schedulepublicationsconstraints "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/constraints/schedule_publications_constraints"
	userstobillingplansconstraints "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/constraints/users_to_billing_plans_constraints"
)

var MigratingConstraintSQLs = []string{
	userstobillingplansconstraints.UserIdBillingPlanIdPartialStatusIndexSQL,
	companyjoinrequestsconstraints.CompanyIdRequesterUserIdPendingIndexSQL,
	schedulepublicationsconstraints.CompanyIdWeekStartIndexSQL,
}
