-- 0000_plan_limitation_seed.up.sql
INSERT INTO "PlanLimitationTable" (
    key,
    updated_at,
    created_at
) VALUES
('Free', NOW(), NOW()),
('Pro', NOW(), NOW()),
('Premium', NOW(), NOW()),
('Ultimate', NOW(), NOW()),
('Enterprise', NOW(), NOW())
ON CONFLICT (key) DO UPDATE SET
    updated_at = NOW();
