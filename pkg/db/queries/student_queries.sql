-- name: GetEligibleStudents :many
SELECT student_table.student_id, usn, name, branch, cgpa, num_active_backlogs, email_id,counsellor_email_id
FROM student_table
JOIN student_job_application_table
ON student_table.student_id = student_job_application_table.student_id
WHERE job_id = $1;

-- name: InsertUser :exec
INSERT INTO student_table(student_id, usn, name, branch, cgpa, num_active_backlogs, email_id,counsellor_email_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);