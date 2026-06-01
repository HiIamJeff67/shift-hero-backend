package constraints

import (
	userstobillingplansconstraints "github.com/your-org/go-start-monolithic-kit/app/models/schemas/constraints/users_to_billing_plans_constraints"
)

var MigratingConstraintSQLs = []string{
	userstobillingplansconstraints.UserIdBillingPlanIdPartialStatusIndexSQL,
}
