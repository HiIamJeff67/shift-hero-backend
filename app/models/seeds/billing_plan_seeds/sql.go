package billingplanseeds

import (
	_ "embed"
)

//go:embed 0000_billing_plan_seed.example.sql
var BillingPlanSeedingExampleDataSQL_0000 string
