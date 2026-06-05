package schedulepublicationsconstraints

import (
	_ "embed"
)

var (
	//go:embed company_id_week_start_idx.sql
	CompanyIdWeekStartIndexSQL string
)
