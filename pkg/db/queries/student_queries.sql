-- name: ApplyForJob :exec
INSERT INTO student_job_application_table(student_id, job_id, applied_on_date)
VALUES($1, $2, NOW());

-- name: GetJobOffers :many
SELECT job_table.job_id, company_name, job_role, job_type, ctc, salary_tier, action, action_date, act_by_date
FROM student_offer_table JOIN job_table 
ON student_offer_table.job_id = job_table.job_id
JOIN company_table
ON job_table.company_id = company_table.company_id
AND student_id = $1;

-- name: GetJobOfferActByDate :one
SELECT act_by_date 
FROM student_offer_table
WHERE student_id = $1 
AND job_id = $2;

-- name: PerformJobOfferAction :exec
UPDATE student_offer_table SET action = $3, action_date = NOW() 
WHERE student_id = $1
AND job_id = $2;

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
  AND (COALESCE(array_length(sqlc.arg(already_applied_job_id)::INT[], 1), 0) = 0 OR job_id <> ANY(sqlc.arg(already_applied_job_id)))
  AND ARRAY(SELECT branch FROM student_table WHERE student_id = $1) && eligible_branches
AND job_table.eligible_batch = (SELECT batch from student_table where student_id = $1);


-- name: GetStudentProfile :one
SELECT name, usn, branch, cgpa, batch, num_active_backlogs, email_id, counsellor_email_id
FROM student_table
WHERE student_id = sqlc.arg(student_id);

-- name: CheckIfAppliedForJobAlready :one
SELECT EXISTS(
    SELECT job_id
    from student_job_application_table
    WHERE job_id = sqlc.arg(job_id)
    AND student_id = sqlc.arg(student_id)
);

-- name: GetAlreadyAppliedJobs :many
SELECT j.job_id, job_role, job_type, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches
FROM student_job_application_table as sj
JOIN job_table as j
ON sj.job_id = j.job_id
WHERE sj.student_id = sqlc.arg(student_id);

-- name: GetAlreadyAppliedJobIds :many
SELECT job_id
FROM student_job_application_table
WHERE student_id = sqlc.arg(student_id);