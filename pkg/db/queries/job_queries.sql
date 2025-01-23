-- name: GetJobs :many
SELECT job_id, job_role, ctc, salary_tier, apply_by_date, cgpa_cutoff, company_name, industry,
       (CASE WHEN 
       cgpa_cutoff <= (SELECT cgpa FROM student_table WHERE student_id = $1) THEN TRUE 
       ELSE FALSE 
       END) AS can_apply
FROM job_table JOIN company_table
ON job_table.company_id = company_table.company_id
WHERE (COALESCE(array_length($2::VARCHAR[], 1), 0) = 0 OR salary_tier = ANY($2))
AND (COALESCE(array_length($3::VARCHAR[], 1), 0) = 0 OR job_role = ANY($3))
AND (COALESCE(array_length($4::VARCHAR[], 1), 0) = 0 OR company_name = ANY($4))
AND NOW() < apply_by_date
AND ARRAY(SELECT branch FROM student_table WHERE student_id = $1) && eligible_branches
AND job_table.eligible_batch = (SELECT batch from student_table where student_id = $1);

-- name: CreateJob :exec
INSERT INTO job_table(company_id, job_role, job_type, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: OfferJob :exec
INSERT INTO student_offer_table(student_id, job_id, action, action_date, act_by_date)
VALUES($1, $2, 'PENDING', NULL, $3);
