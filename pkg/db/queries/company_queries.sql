-- name: CreateCompany :exec
INSERT INTO company_table (company_name, poc_name, poc_phno, industry, username, password)
VALUES ($1, $2, $3, $4, $5, $6);
