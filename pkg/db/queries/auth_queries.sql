-- name: AuthenticateStudent :one
SELECT student_id from student_table where usn = $1 and password = $2;

-- name: AuthenticateCompany :one
SELECT company_id from company_table where username = $1 and password = $2;

