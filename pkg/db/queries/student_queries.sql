-- name: GetEligibleStudents :many
SELECT student_table.student_id, usn, name, branch, cgpa, num_active_backlogs, email_id
FROM student_job_application_table
JOIN student_table
ON student_table.student_id = student_job_application_table.student_id
WHERE job_id = $1;


-- name: RegisterForJob :exec
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


