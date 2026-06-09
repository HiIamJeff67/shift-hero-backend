SELECT
    user_account.user_id,
    user_account.ai_monthly_usage_count AS monthly_usage_count,
    user_account.ai_usage_period_start AS period_start,
    plan_limitation.ai_monthly_generation_limit AS monthly_limit
FROM "UserAccountTable" AS user_account
JOIN "UserTable" AS app_user
    ON app_user.id = user_account.user_id
JOIN "PlanLimitationTable" AS plan_limitation
    ON plan_limitation.key = app_user.plan
WHERE user_account.user_id = @user_id
FOR UPDATE OF user_account;
