ALTER TABLE "UserAccountTable"
DROP CONSTRAINT IF EXISTS user_account_ai_monthly_usage_count_nonnegative;

-- ============================== SQL Separator ==============================

ALTER TABLE "UserAccountTable"
ADD CONSTRAINT user_account_ai_monthly_usage_count_nonnegative
CHECK (ai_monthly_usage_count >= 0);

-- ============================== SQL Separator ==============================

ALTER TABLE "PlanLimitationTable"
DROP CONSTRAINT IF EXISTS plan_limitation_ai_monthly_generation_limit_nonnegative;

-- ============================== SQL Separator ==============================

ALTER TABLE "PlanLimitationTable"
ADD CONSTRAINT plan_limitation_ai_monthly_generation_limit_nonnegative
CHECK (ai_monthly_generation_limit >= 0);
