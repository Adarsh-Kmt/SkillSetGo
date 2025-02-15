-- name: CreateJob :exec
INSERT INTO job_table(company_id, job_role, job_type, job_description, ctc, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: OfferJob :exec
INSERT INTO student_offer_table(student_id, job_id, action, action_date, act_by_date)
VALUES($1, $2, 'PENDING', NULL, $3);

-- name: CheckIfOfferedAlready :one
SELECT EXISTS(
    SELECT student_id
    FROM student_offer_table
    WHERE student_id = sqlc.arg(student_id)
    AND job_id = sqlc.arg(job_id)
);

-- name: GetPublishedJobs :many
SELECT job_id, job_role, job_type, ctc, job_description, salary_tier, apply_by_date, cgpa_cutoff, eligible_batch, eligible_branches
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

-- name: GetOfferStatus :many
SELECT s.student_id, usn, name, branch, j.job_id, j.job_role, action, action_date, act_by_date
FROM student_offer_table as so
JOIN student_table as s
ON s.student_id = so.student_id
JOIN job_table as j
ON so.job_id = j.job_id
WHERE so.job_id = sqlc.arg(job_id);

-- name: CheckIfInterviewExists :one
SELECT EXISTS(
    SELECT student_id
    FROM student_job_interview_table
    WHERE student_id = sqlc.arg(student_id)
    AND job_id = sqlc.arg(job_id)
    AND interview_round = sqlc.arg(interview_round)
);

-- name: UpdateInterviewResult :exec
UPDATE student_job_interview_table SET result = sqlc.arg(result)
WHERE student_id = sqlc.arg(student_id)
AND job_id = sqlc.arg(job_id)
AND interview_round = sqlc.arg(interview_round);

-- name: GetInterviewsScheduledForStudent :many
SELECT sj.job_id, job_role, job_type, ctc, company_name, venue, interview_date
FROM student_job_interview_table as sj
JOIN job_table as j
ON j.job_id = sj.job_id
JOIN company_table as c
ON j.company_id = c.company_id
WHERE sj.student_id = sqlc.arg(student_id);

-- name: GetInterviewsScheduledByCompany :many
SELECT s.name, s.usn, s.student_id, s.cgpa, venue, interview_date
FROM student_job_interview_table as sj
JOIN student_table as s
ON sj.student_id = s.student_id
WHERE job_id = sqlc.arg(job_id);

-- name: CheckIfVenueBeingUsedAtParticularTime :one
SELECT EXISTS(
    SELECT student_id
    FROM student_job_interview_table
    WHERE venue = sqlc.arg(venue)
    AND ABS(EXTRACT(EPOCH FROM (sqlc.arg(date_of_interview_to_be_scheduled) - interview_date))) <= 1800
);

-- name: ScheduleInterview :exec
INSERT INTO student_job_interview_table(student_id, job_id, interview_date, venue)
VALUES($1, $2, $3, $4);


-- name: GetPlacementStats :many
SELECT s.branch, AVG(j.ctc) as mean_lpa, COUNT(so.student_id) as number_of_offers
FROM student_offer_table as so
JOIN job_table as j
ON so.job_id = j.job_id
JOIN student_table as s
ON so.student_id = s.student_id
GROUP BY branch;