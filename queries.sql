-- name: InsertUser :exec
INSERT INTO users (name, branch, cgpa, active_backlogs, email_id, usn, counsellor_name)
VALUES ($1, $2, $3, $4, $5, $6, $7);