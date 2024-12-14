-- name: GetJobs :many
SELECT job_id, job_role, ctc, salary_tier, apply_by_date, cgpa_cutoff, company_name, industry
FROM job_table JOIN company_table
on job_table.company_id = company_table.company_id
where (COALESCE(array_length($2::VARCHAR[], 1), 0) = 0 OR salary_tier = ANY($2))
and NOW() < apply_by_date
and cgpa_cutoff <= (SELECT cgpa from student_table where student_id = $1);

-- name: CreateJob :exec
INSERT INTO job_table(company_id, job_role, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_branches)
VALUES($1, $2, $3, $4, $5, $6, $7);