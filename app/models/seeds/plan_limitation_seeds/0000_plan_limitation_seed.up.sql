-- 0000_plan_limitation_seed.up.sql
INSERT INTO "PlanLimitationTable" (
    key,
    ai_monthly_generation_limit,
    updated_at,
    created_at
) VALUES
('Free', 5, NOW(), NOW()),
('Pro', 30, NOW(), NOW()),
('Premium', 100, NOW(), NOW()),
('Ultimate', 300, NOW(), NOW()),
('Enterprise', 1000, NOW(), NOW())
ON CONFLICT (key) DO UPDATE SET
    ai_monthly_generation_limit = EXCLUDED.ai_monthly_generation_limit,
    updated_at = NOW();
