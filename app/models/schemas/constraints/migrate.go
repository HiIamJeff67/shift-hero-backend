package constraints

import (
	userstobillingplansconstraints "github.com/HiIamJeff67/shift-hero-backend/app/models/schemas/constraints/users_to_billing_plans_constraints"
)

var MigratingConstraintSQLs = []string{
	userstobillingplansconstraints.UserIdBillingPlanIdPartialStatusIndexSQL,
}
