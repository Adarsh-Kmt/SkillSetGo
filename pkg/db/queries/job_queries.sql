-- name: GetJobs :many
SELECT job_id, job_role, ctc, salary_tier, apply_by_date, cgpa_cutoff, company_name, industry,
       (CASE WHEN cgpa_cutoff <= (SELECT cgpa FROM student_table WHERE student_id = $1) THEN TRUE ELSE FALSE END) AS can_apply
FROM job_table JOIN company_table
ON job_table.company_id = company_table.company_id
-- apply salary tier filter
WHERE (COALESCE(array_length($2::VARCHAR[], 1), 0) = 0 OR salary_tier = ANY($2))
-- apply job role filter
AND (COALESCE(array_length($3::VARCHAR[], 1), 0) = 0 OR job_role = ANY($3))
-- apply company filter
AND (COALESCE(array_length($4::VARCHAR[], 1), 0) = 0 OR company_name = ANY($4))
-- filter out companies whose application date expired.
AND NOW() < apply_by_date
-- filter out companies for which the student's branch is not eligible
AND ARRAY(SELECT branch FROM student_table WHERE student_id = $1) && eligible_branches;
-- and cgpa_cutoff <= (SELECT cgpa from student_table where student_id = $1);

-- name: CreateJob :exec
INSERT INTO job_table(company_id, job_role, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_branches)
VALUES($1, $2, $3, $4, $5, $6, $7);