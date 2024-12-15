-- name: GetEligibleStudents :many
SELECT student_id, usn, name, branch, cgpa, num_active_backlogs, email_id
FROM student_job_application_table
JOIN student_table
ON student_table.student_id = student_job_application_table.student_id
WHERE job_id = $1;
