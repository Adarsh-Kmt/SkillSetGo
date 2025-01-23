-- name: AuthenticateStudent :one
SELECT student_id from student_table where usn = $1 and password = $2;

-- name: AuthenticateCompany :one
SELECT company_id from company_table where username = $1 and password = $2;

-- name: CreateCompany :exec
INSERT INTO company_table (company_name, poc_name, poc_phno, industry, username, password)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: InsertUser :exec
INSERT INTO student_table(usn, name, password, branch, batch,cgpa, num_active_backlogs,email_id,counsellor_email_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);