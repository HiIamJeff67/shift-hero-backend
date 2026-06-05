DROP INDEX IF EXISTS "company_join_requests_idx_company_id_requester_user_id_pending";

-- ============================== SQL Separator ==============================

CREATE UNIQUE INDEX company_join_requests_idx_company_id_requester_user_id_pending
ON "CompanyJoinRequestsTable" (company_id, requester_user_id)
WHERE status = 'Pending';
