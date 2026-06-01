package seeds

import (
	_ "embed"

	billingplanseeds "github.com/your-org/go-start-monolithic-kit/app/models/seeds/billing_plan_seeds"
	planlimitationseeds "github.com/your-org/go-start-monolithic-kit/app/models/seeds/plan_limitation_seeds"
)

var SeedingDefaultDataSQLs = []string{
	planlimitationseeds.PlanLimitationSeedingDefaultDataSQL_0000_UP,
	billingplanseeds.BillingPlanSeedingExampleDataSQL_0000,
}
