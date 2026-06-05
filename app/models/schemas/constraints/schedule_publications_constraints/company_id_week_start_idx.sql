DROP INDEX IF EXISTS "schedule_publications_idx_company_id_week_start";

-- ============================== SQL Separator ==============================

CREATE UNIQUE INDEX schedule_publications_idx_company_id_week_start
ON "SchedulePublicationsTable" (company_id, week_start);
