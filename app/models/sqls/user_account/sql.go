package useraccountsqls

import (
	_ "embed"
)

//go:embed get_ai_usage_quota_by_user_id_for_update.sql
var GetAIUsageQuotaByUserIdForUpdateSQL string
