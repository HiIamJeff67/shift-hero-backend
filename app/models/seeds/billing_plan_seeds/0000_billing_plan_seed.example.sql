-- 0000_billing_plan_seed.example.sql
-- IMPORTANT:
-- This is a template seed example only.
-- Replace ids, product_ids, and names with your real billing provider resources.
-- The "name" values must stay aligned with app/models/schemas/enums/billing_plan_name_enum.go.

INSERT INTO "BillingPlanTable" (
    id,
    product_id,
    name,
    status,
    interval_unit,
    price,
    currency_code,
    updated_at,
    created_at
) VALUES
('BP-TEMPLATE-MONTHLY-FREE',       'PROD-TEMPLATE-CORE', 'Template Monthly Free Plan',       'ACTIVE', 'MONTH', 0.00,   'USD', NOW(), NOW()),
('BP-TEMPLATE-MONTHLY-PRO',        'PROD-TEMPLATE-CORE', 'Template Monthly Pro Plan',        'ACTIVE', 'MONTH', 4.99,   'USD', NOW(), NOW()),
('BP-TEMPLATE-YEARLY-PRO',         'PROD-TEMPLATE-CORE', 'Template Yearly Pro Plan',         'ACTIVE', 'YEAR',  49.99,  'USD', NOW(), NOW()),
('BP-TEMPLATE-MONTHLY-PREMIUM',    'PROD-TEMPLATE-CORE', 'Template Monthly Premium Plan',    'ACTIVE', 'MONTH', 9.99,   'USD', NOW(), NOW()),
('BP-TEMPLATE-YEARLY-PREMIUM',     'PROD-TEMPLATE-CORE', 'Template Yearly Premium Plan',     'ACTIVE', 'YEAR',  99.99,  'USD', NOW(), NOW()),
('BP-TEMPLATE-MONTHLY-ULTIMATE',   'PROD-TEMPLATE-CORE', 'Template Monthly Ultimate Plan',   'ACTIVE', 'MONTH', 19.99,  'USD', NOW(), NOW()),
('BP-TEMPLATE-YEARLY-ULTIMATE',    'PROD-TEMPLATE-CORE', 'Template Yearly Ultimate Plan',    'ACTIVE', 'YEAR',  199.99, 'USD', NOW(), NOW()),
('BP-TEMPLATE-MONTHLY-ENTERPRISE', 'PROD-TEMPLATE-CORE', 'Template Monthly Enterprise Plan', 'ACTIVE', 'MONTH', 49.99,  'USD', NOW(), NOW()),
('BP-TEMPLATE-YEARLY-ENTERPRISE',  'PROD-TEMPLATE-CORE', 'Template Yearly Enterprise Plan',  'ACTIVE', 'YEAR',  499.99, 'USD', NOW(), NOW())
ON CONFLICT (name) DO UPDATE SET
    id = EXCLUDED.id,
    product_id = EXCLUDED.product_id,
    status = EXCLUDED.status,
    interval_unit = EXCLUDED.interval_unit,
    price = EXCLUDED.price,
    currency_code = EXCLUDED.currency_code,
    updated_at = NOW();
