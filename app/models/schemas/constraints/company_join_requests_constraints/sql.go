package companyjoinrequestsconstraints

import (
	_ "embed"
)

var (
	//go:embed company_id_requester_user_id_pending_idx.sql
	CompanyIdRequesterUserIdPendingIndexSQL string
)
