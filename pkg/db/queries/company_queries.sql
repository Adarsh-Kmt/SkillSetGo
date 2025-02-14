-- name: CreateJob :exec
INSERT INTO job_table(company_id, job_role, job_type, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: OfferJob :exec
INSERT INTO student_offer_table(student_id, job_id, action, action_date, act_by_date)
VALUES($1, $2, 'PENDING', NULL, $3);

-- name: GetPublishedJobs :many
SELECT job_id, job_role, job_type, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches
FROM job_table
WHERE company_id = sqlc.arg(company_id);

-- name: GetEligibleStudents :many
SELECT student_table.student_id, usn, name, branch, cgpa, num_active_backlogs, email_id
FROM student_job_application_table
         JOIN student_table
              ON student_table.student_id = student_job_application_table.student_id
WHERE job_id = $1;

-- name: GetJobApplicants :many
SELECT s.student_id, name, usn, branch, cgpa, batch, num_active_backlogs, email_id, counsellor_email_id
FROM student_table as s
JOIN student_job_application_table as sj
ON s.student_id = sj.student_id
WHERE job_id = sqlc.arg(job_id);

-- name: CheckIfCompanyCreatedJob :one
SELECT EXISTS(
    SELECT company_id
    FROM job_table
    WHERE company_id = sqlc.arg(company_id)
    AND job_id = sqlc.arg(job_id)
);